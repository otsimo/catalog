package catalog

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sync"

	phttp "github.com/coreos/go-oidc/http"
	"github.com/coreos/go-oidc/jose"
	"github.com/coreos/go-oidc/key"
	"github.com/coreos/go-oidc/oidc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	keySyncWindow      = 5 * time.Second
	UserGroupAdmin     = "otsimo.com/admin"
	UserGroupOwner     = "otsimo.com/owner"
	UserGroupDeveloper = "otsimo.com/developer"
)

func NewOIDCClient(id, secret, discovery string) (*Client, error) {
	var cfg oidc.ProviderConfig
	var err error
	for {
		cfg, err = oidc.FetchProviderConfig(http.DefaultClient, discovery)
		if err == nil {
			break
		}
		sleep := 1 * time.Second
		fmt.Printf("Failed fetching provider config, trying again in %v: %v\n", sleep, err)
		time.Sleep(sleep)
	}

	c := Client{
		providerConfig: newProviderConfigRepo(cfg),
		httpClient:     http.DefaultClient,
	}

	return &c, nil
}

func getJWTToken(ctx context.Context) (jose.JWT, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return jose.JWT{}, fmt.Errorf("missing metadata")
	}
	var auth []string
	auth, ok = md["authorization"]
	if !ok || len(auth) == 0 {
		auth, ok = md["Authorization"]
		if !ok || len(auth) == 0 {
			return jose.JWT{}, errors.New("missing authorization header")
		}
	}
	ah := auth[0]
	if len(ah) <= 6 || strings.ToUpper(ah[0:6]) != "BEARER" {
		return jose.JWT{}, errors.New("should be a bearer token")
	}
	val := ah[7:]
	if len(val) == 0 {
		return jose.JWT{}, errors.New("bearer token is empty")
	}
	return jose.ParseJWT(val)
}

func hasOneOf(userGroups, wantedGroups []string) bool {
	if len(wantedGroups) == 0 {
		return true
	}
	for _, ug := range userGroups {
		for _, wg := range wantedGroups {
			if ug == wg {
				return true
			}
		}
	}
	return false
}

func (s *catalogGrpcServer) authToken(jwt jose.JWT, accessGroups []string) (string, string, error) {
	claims, err := jwt.Claims()
	if err != nil {
		return "", "", fmt.Errorf("auth.go: failed to get claims %v", err)
	}

	aud, ok, err := claims.StringClaim("aud")
	if err != nil || !ok || aud == "" {
		return "", "", fmt.Errorf("auth.go: failed to parse 'aud' claim: %v", err)
	}

	err = s.server.Oidc.VerifyJWT(jwt, aud)
	if err != nil {
		return "", "", fmt.Errorf("auth.go: Failed to verify signature: %v", err)
	}

	sub, ok, err := claims.StringClaim("sub")
	if err != nil {
		return "", "", fmt.Errorf("auth.go: failed to parse 'sub' claim: %v", err)
	}
	if !ok || sub == "" {
		return "", "", fmt.Errorf("auth.go: missing required 'sub' claim")
	}

	email, ok, err := claims.StringClaim("email")
	if err != nil || !ok || email == "" {
		return "", "", fmt.Errorf("auth.go: failed to parse 'email' claim: %v", err)
	}
	groups, ok, err := claims.StringsClaim("groups")
	if err != nil {
		return "", "", fmt.Errorf("auth.go: failed to parse claim: %v", err)
	}
	if !hasOneOf(groups, accessGroups) {
		return "", "", fmt.Errorf("auth.go: user does not have required permission")
	}
	return sub, email, nil
}

type Client struct {
	httpClient     phttp.Client
	providerConfig *providerConfigRepo
	keySet         key.PublicKeySet
	providerSyncer *oidc.ProviderConfigSyncer

	keySetSyncMutex sync.RWMutex
	lastKeySetSync  time.Time
}

func (c *Client) Healthy() error {
	now := time.Now().UTC()

	cfg := c.providerConfig.Get()

	if cfg.Empty() {
		return errors.New("oidc client provider config empty")
	}

	if !cfg.ExpiresAt.IsZero() && cfg.ExpiresAt.Before(now) {
		return errors.New("oidc client provider config expired")
	}

	return nil
}

// SyncProviderConfig starts the provider config syncer
func (c *Client) SyncProviderConfig(discoveryURL string) chan struct{} {
	r := oidc.NewHTTPProviderConfigGetter(c.httpClient, discoveryURL)
	s := oidc.NewProviderConfigSyncer(r, c.providerConfig)
	stop := s.Run()
	s.WaitUntilInitialSync()
	return stop
}

func (c *Client) maybeSyncKeys() error {
	tooSoon := func() bool {
		return time.Now().UTC().Before(c.lastKeySetSync.Add(keySyncWindow))
	}

	// ignore request to sync keys if a sync operation has been
	// attempted too recently
	if tooSoon() {
		return nil
	}

	c.keySetSyncMutex.Lock()
	defer c.keySetSyncMutex.Unlock()

	// check again, as another goroutine may have been holding
	// the lock while updating the keys
	if tooSoon() {
		return nil
	}

	cfg := c.providerConfig.Get()
	r := oidc.NewRemotePublicKeyRepo(c.httpClient, cfg.KeysEndpoint.String())
	w := &clientKeyRepo{client: c}
	_, err := key.Sync(r, w)
	c.lastKeySetSync = time.Now().UTC()

	return err
}

type clientKeyRepo struct {
	client *Client
}

func (r *clientKeyRepo) Set(ks key.KeySet) error {
	pks, ok := ks.(*key.PublicKeySet)
	if !ok {
		return errors.New("unable to cast to PublicKey")
	}
	r.client.keySet = *pks
	return nil
}

func (c *Client) VerifyJWT(jwt jose.JWT, clientID string) error {
	var keysFunc func() []key.PublicKey
	if kID, ok := jwt.KeyID(); ok {
		keysFunc = c.keysFuncWithID(kID)
	} else {
		keysFunc = c.keysFuncAll()
	}

	v := oidc.NewJWTVerifier(
		c.providerConfig.Get().Issuer.String(),
		clientID,
		c.maybeSyncKeys, keysFunc)

	return v.Verify(jwt)
}

// keysFuncWithID returns a function that retrieves at most unexpired
// public key from the Client that matches the provided ID
func (c *Client) keysFuncWithID(kID string) func() []key.PublicKey {
	return func() []key.PublicKey {
		c.keySetSyncMutex.RLock()
		defer c.keySetSyncMutex.RUnlock()

		if c.keySet.ExpiresAt().Before(time.Now()) {
			return []key.PublicKey{}
		}

		k := c.keySet.Key(kID)
		if k == nil {
			return []key.PublicKey{}
		}

		return []key.PublicKey{*k}
	}
}

// keysFuncAll returns a function that retrieves all unexpired public
// keys from the Client
func (c *Client) keysFuncAll() func() []key.PublicKey {
	return func() []key.PublicKey {
		c.keySetSyncMutex.RLock()
		defer c.keySetSyncMutex.RUnlock()

		if c.keySet.ExpiresAt().Before(time.Now()) {
			return []key.PublicKey{}
		}

		return c.keySet.Keys()
	}
}

type providerConfigRepo struct {
	mu     sync.RWMutex
	config oidc.ProviderConfig // do not access directly, use Get()
}

func newProviderConfigRepo(pc oidc.ProviderConfig) *providerConfigRepo {
	return &providerConfigRepo{sync.RWMutex{}, pc}
}

func (r *providerConfigRepo) Set(cfg oidc.ProviderConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.config = cfg
	return nil
}

func (r *providerConfigRepo) Get() oidc.ProviderConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.config
}
