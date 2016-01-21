package main

import (
	"fmt"

	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/bgentry/speakeasy"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
)

// oauthAccess supplies credentials from a given token.
type oauthAccess struct {
	Config     *Config
	RequireTLS bool
}

// NewOauthAccess constructs the credentials using a given token.
func NewOauthAccess(cnf *Config) oauthAccess {
	return oauthAccess{Config: cnf, RequireTLS: true}
}

func (oa *oauthAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	token := oa.Config.Token(accountHost)
	if len(token) == 0 {
		logrus.Fatalln("user not logged in")
	}
	return map[string]string{
		"Authorization": "Bearer" + " " + token,
	}, nil
}

func (oa *oauthAccess) RequireTransportSecurity() bool {
	return oa.RequireTLS
}

func login(ctx *cli.Context) {
	email := ctx.String("email")
	password := ctx.String("password")
	if len(email) == 0 {
		for {
			fmt.Printf("your email: ")
			if c, err := fmt.Scanf("%s", &email); err == nil {
				break
			} else if c == 0 && err.Error() == "unexpected newline" {
				break
			}
			fmt.Println("invalid input")
		}
	}
	var err error
	if len(password) == 0 {
		for {
			password, err = speakeasy.Ask("your password: ")
			if len(password) == 0 {
				fmt.Println("invalid input")
				continue
			}
			if err == nil {
				break
			}
		}
	}
	val := url.Values{}
	val.Add("username", email)
	val.Add("password", password)
	val.Add("grant_type", "password")
	resp, err := http.PostForm(accountHost+"/login", val)
	if err != nil {
		logrus.Fatalln("error while trying to login", err)
	}
	dec := json.NewDecoder(resp.Body)
	asd := struct {
		Error        string `json:"error,omitempty"`
		AccessToken  string `json:"access_token,omitempty"`
		TokenType    string `json:"token_type,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}{
		AccessToken:  "",
		TokenType:    "",
		RefreshToken: "",
	}
	err = dec.Decode(&asd)
	if err != nil {
		logrus.Fatalln("unable to decode response body", err)
	}
	if len(asd.Error) > 0 {
		logrus.Fatalln("unable to login", asd.Error)
	}
	cnf := config()
	cnf.Auths[accountHost] = AuthConfig{Email: email, Token: asd.AccessToken}
	err = cnf.Save()
	if err != nil {
		logrus.Fatalln("unable to save config", err)
	}
	logrus.Infoln("logged in to user=", email)
}
