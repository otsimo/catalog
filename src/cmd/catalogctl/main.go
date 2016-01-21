package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/otsimo/api/apipb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	Version string
	conn      *grpc.ClientConn
	client apipb.CatalogServiceClient
	cafile string
	remoteUrl string = "127.0.0.1:18857"
)

func connect() {
	var opts []grpc.DialOption

	if len(cafile) > 0 {
		auth, err := credentials.NewClientTLSFromFile(cafile, "")
		if err != nil {
			panic(err)
		} else {
			opts = append(opts, grpc.WithTransportCredentials(auth))
		}
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
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

func initialize(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	cafile = ctx.String("tls-ca-file")
	remoteUrl = ctx.String("url")

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
	//todo add auth header
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
	log.Infof("Current Catalog:\n%+v", res)
}

func main() {
	app := cli.NewApp()
	app.Name = "catalogctl"
	app.Version = Version
	app.Usage = "Otsimo Catalog Service Client"
	app.Author = "Sercan Değirmenci <sercan@otsimo.com>"
	var flags []cli.Flag

	flags = []cli.Flag{
		cli.StringFlag{Name: "url, u", Value: "127.0.0.1:18857", Usage: "remote server url"},
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
	}
	app.Run(os.Args)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
