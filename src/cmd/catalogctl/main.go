package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/otsimo/api/apipb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"fmt"
	"encoding/json"
)

var (
	Version string
	conn        *grpc.ClientConn
	client apipb.CatalogServiceClient
	cafile string
	remoteUrl string = "127.0.0.1:18857"
	accountHost string = "http://127.0.0.1:18856"
)

func connect() {
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
	conn, err := grpc.Dial(remoteUrl, opts...)

	if err != nil {
		log.Fatalf("main.go: Error while connection to catalog service %v\n", err)
	}
	client = apipb.NewCatalogServiceClient(conn)
}

func closeConn() {
	if conn != nil {
		conn.Close()
		conn = nil
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
	r, err := cf.Request()
	if err != nil {
		log.Fatalln("error while creating push catalog request, error=", err)
	}
	connect()
	defer closeConn()
	resp, err := client.Push(context.Background(), r)
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
	log.Info("catalog is valid")
}

func current(ctx *cli.Context) {
	connect()
	defer closeConn()
	res, err := client.Pull(context.Background(), &apipb.CatalogPullRequest{})
	if err != nil {
		log.Fatalln(err)
	}
	b, _ := json.MarshalIndent(res,"","  ")
	fmt.Printf("Current Catalog:\n%s", string(b))
}

func approve(ctx *cli.Context) {
	connect()
	defer closeConn()
	if !ctx.Args().Present() {
		log.Fatalln("enter a valid catalog title")
	}
	title := ctx.Args().First()
	_, err := client.Approve(context.Background(), &apipb.CatalogApproveRequest{Title: title})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Catalog '%s' approved\n", title)
}

func list(ctx *cli.Context) {
	connect()
	defer closeConn()
	query := &apipb.CatalogListRequest{
		HideExpired : ctx.Bool("hide-expired"),
		Limit: int32(ctx.Int("limit")),
	}
	stat := ctx.String("status")
	if stat == "draft" {
		query.Status = apipb.CatalogListRequest_ONLY_DRAFT
	}else if stat == "approved" {
		query.Status = apipb.CatalogListRequest_ONLY_APPROVED
	}else if stat == "both" {
		query.Status = apipb.CatalogListRequest_BOTH
	}
	res, err := client.List(context.Background(), query)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Found %d catalog(s)\n", len(res.Catalogs))
	fmt.Println("index\ttitle\tstatus")
	for i, v := range res.Catalogs {
		fmt.Printf("%d\t%s\t%s\n", (i + 1), v.Title, apipb.CatalogStatus_name[int32(v.Status)])
	}
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
		cli.StringFlag{Name: "tls-ca-file", Value: "", Usage: "the server's certificate file for TLS connection"},
		cli.BoolFlag{Name: "debug, d", Usage: "enable verbose log"},
	}

	app.Flags = flags
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
				cli.StringFlag{Name:"status", Value:"both", Usage:"catalog status: both, draft, approved"},
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
