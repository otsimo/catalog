package main

import (
	"os"

	"encoding/json"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/otsimo/api/apipb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	Version string
	catalogConn        *grpc.ClientConn
	registryConn        *grpc.ClientConn
	catalogClient apipb.CatalogServiceClient
	registryClient apipb.RegistryServiceClient
	cafile string
	remoteUrl string = "catalog.otsimo.com"
	registryUrl string = "registry.otsimo.com"
	accountHost string = "https://accounts.otsimo.com"
)

func connect(connect2registry bool) {
	var opts []grpc.DialOption
	jwtCreds := NewOauthAccess(config())
	if len(cafile) > 0 {
		auth, err := credentials.NewClientTLSFromFile(cafile, "")
		if err != nil {
			panic(err)
		} else {
			opts = append(opts, grpc.WithTransportCredentials(auth))
		}
	} else {
		jwtCreds.RequireTLS = false
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithPerRPCCredentials(&jwtCreds))
	var err error
	catalogConn, err = grpc.Dial(remoteUrl, opts...)
	if err != nil {
		log.Fatalf("main.go: Error while connection to catalog service %v\n", err)
	}
	catalogClient = apipb.NewCatalogServiceClient(catalogConn)

	if connect2registry {
		registryConn, err = grpc.Dial(registryUrl, opts...)
		if err != nil {
			log.Fatalf("main.go: Error while connection to registry service %v\n", err)
		}
		registryClient = apipb.NewRegistryServiceClient(registryConn)
	}
}

func closeConn() {
	if catalogConn != nil {
		catalogConn.Close()
		catalogConn = nil
	}
	if registryConn != nil {
		registryConn.Close()
		registryConn = nil
	}
}

func config() *Config {
	c, err := NewConfig()
	if err != nil {
		log.Fatalf("Unable create config, error=%+v", err)
	}
	return c
}

func initialize(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	cafile = ctx.String("tls-ca-file")
	remoteUrl = ctx.String("url")
	accountHost = ctx.String("auth")
	registryUrl = ctx.String("registry")
	return nil
}

func push(ctx *cli.Context) {
	if !ctx.Args().Present() {
		log.Fatalln("enter a valid catalog file path")
	}
	cf, err := readCatalogFile(ctx.Args().First())
	if err != nil {
		log.Fatalln("error while reading catalog file, error=", err)
	}
	connect(true)
	defer closeConn()

	r, err := cf.Request()
	if err != nil {
		log.Fatalln("error while creating push catalog request, error=", err)
	}

	resp, err := catalogClient.Push(context.Background(), r)
	if err != nil {
		log.Fatalln("response returned with error=", err)
	}
	log.Infof("pushed, result: %+v\n", resp)
}

func validate(ctx *cli.Context) {
	if !ctx.Args().Present() {
		log.Fatalln("enter a valid catalog file path")
	}
	_, err := readCatalogFile(ctx.Args().First())
	if err != nil {
		log.Fatalln("error while reading catalog file, error:", err)
	}
	fmt.Println("catalog is valid")
}

func current(ctx *cli.Context) {
	connect(false)
	defer closeConn()
	res, err := catalogClient.Pull(context.Background(), &apipb.CatalogPullRequest{})
	if err != nil {
		log.Fatalln(err)
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("Current Catalog:\n%s", string(b))
}

func approve(ctx *cli.Context) {
	connect(false)
	defer closeConn()
	if !ctx.Args().Present() {
		log.Fatalln("enter a valid catalog title")
	}
	title := ctx.Args().First()
	_, err := catalogClient.Approve(context.Background(), &apipb.CatalogApproveRequest{Title: title})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Catalog '%s' approved\n", title)
}

func list(ctx *cli.Context) {
	connect(false)
	defer closeConn()
	query := &apipb.CatalogListRequest{
		HideExpired: ctx.Bool("hide-expired"),
		Limit:       int32(ctx.Int("limit")),
	}
	stat := ctx.String("status")
	if stat == "draft" {
		query.Status = apipb.CatalogListRequest_ONLY_DRAFT
	} else if stat == "approved" {
		query.Status = apipb.CatalogListRequest_ONLY_APPROVED
	} else if stat == "both" {
		query.Status = apipb.CatalogListRequest_BOTH
	}
	res, err := catalogClient.List(context.Background(), query)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Found %d catalog(s)\n", len(res.Catalogs))
	fmt.Println("index\ttitle\tstatus")
	for i, v := range res.Catalogs {
		fmt.Printf("%d\t%s\t%s\n", (i + 1), v.Title, apipb.CatalogStatus_name[int32(v.Status)])
	}
}
func withEnvs(prefix string, flags []cli.Flag) []cli.Flag {
	var flgs []cli.Flag
	for _, f := range flags {
		env := ""
		spr := strings.Split(f.GetName(), ",")
		env = prefix + "_" + strings.ToUpper(strings.Replace(spr[0], "-", "_", -1))
		switch v := f.(type) {
		case cli.IntFlag:
			flgs = append(flgs, cli.IntFlag{Name: v.Name, Value: v.Value, Usage: v.Usage, EnvVar: env})
		case cli.StringFlag:
			flgs = append(flgs, cli.StringFlag{Name: v.Name, Value: v.Value, Usage: v.Usage, EnvVar: env})
		case cli.BoolFlag:
			flgs = append(flgs, cli.BoolFlag{Name: v.Name, Usage: v.Usage, EnvVar: env})
		default:
			fmt.Println("unknown")
		}
	}
	return flgs
}

func main() {
	app := cli.NewApp()
	app.Name = "catalogctl"
	app.Version = Version
	app.Usage = "Otsimo Catalog Service Client"
	app.Author = "Sercan Değirmenci <sercan@otsimo.com>"
	var flags []cli.Flag

	flags = []cli.Flag{
		cli.StringFlag{Name: "url, u", Value: remoteUrl, Usage: "remote server url"},
		cli.StringFlag{Name: "auth", Value: accountHost, Usage: "otsimo accounts url"},
		cli.StringFlag{Name: "registry", Value: registryUrl, Usage: "otsimo regitry service url"},
		cli.StringFlag{Name: "tls-ca-file", Value: "", Usage: "the server's certificate file for TLS connection"},
		cli.BoolFlag{Name: "debug, d", Usage: "enable verbose log"},
	}

	app.Flags = withEnvs("CATALOGCTL", flags)
	app.Before = initialize
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:   "push",
			Usage:  "push catalog",
			Action: push,
		},
		{
			Name:   "login",
			Usage:  "login otsimo accounts",
			Action: login,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "email", Value: ""},
				cli.StringFlag{Name: "password", Value: ""},
			},
		},
		{
			Name:   "validate",
			Usage:  "validate catalog file",
			Action: validate,
		},
		{
			Name:   "current",
			Usage:  "get current accessible catalog",
			Action: current,
		},
		{
			Name:   "approve",
			Usage:  "approve a catalog",
			Action: approve,
		},
		{
			Name:   "list",
			Usage:  "list catalogs",
			Action: list,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "hide-expired"},
				cli.IntFlag{Name: "limit", Value: 100},
				cli.StringFlag{Name: "status", Value: "both", Usage: "catalog status: both, draft, approved"},
			},
		},
	}
	app.Run(os.Args)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
