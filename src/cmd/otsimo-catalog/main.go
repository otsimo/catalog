package main

import (
	"catalog"
	"fmt"
	"os"
	"storage"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var Version string
var config *catalog.Config = catalog.NewConfig()

const (
	EnvDebugName    = "OTSIMO_CATALOG_DEBUG"
	EnvGrpcPortName = "OTSIMO_CATALOG_GRPC_PORT"
	EnvStorageName  = "OTSIMO_CATALOG_STORAGE"
)

func RunAction(c *cli.Context) {
	config.Debug = c.Bool("debug")
	config.GrpcPort = c.Int("grpc-port")
	config.TlsCertFile = c.String("tls-cert-file")
	config.TlsKeyFile = c.String("tls-key-file")

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}
	//find driver name
	sname := c.String("storage")
	if sname == "" || sname == "none" {
		log.Errorln("main.go: storage flag is missing or it cannot be 'none'")
		cli.ShowAppHelp(c)
		return
	}

	//get driver
	driver := storage.GetDriver(sname)
	if driver == nil {
		log.Fatalf("main.go: storage driver '%s' not found\n", sname)
	}

	//load storage driver
	s, err := driver.New(c)
	if err != nil {
		log.Fatal("main.go: error while creating new storage driver:", err, s)
	}

	server := catalog.NewServer(config, s)
	server.ListenGRPC()
}

func main() {
	fmt.Println("Otsimo Catalog")

	app := cli.NewApp()
	app.Name = "otsimo-catalog"
	app.Version = Version
	app.Usage = "Otsimo Catalog Service"
	app.Author = "Sercan Değirmenci <sercan@otsimo.com>"
	dnames := storage.GetDriverNames()
	var flags []cli.Flag

	flags = []cli.Flag{
		cli.IntFlag{Name: "grpc-port", Value: catalog.DefaultGrpcPort, Usage: "grpc server port", EnvVar: EnvGrpcPortName},
		cli.StringFlag{Name: "storage, s", Value: "none", Usage: fmt.Sprintf("the storage driver. Available drivers: %s", strings.Join(dnames, ", ")), EnvVar: EnvStorageName},
		cli.StringFlag{Name: "tls-cert-file", Value: "", Usage: "the server's certificate file for TLS connection"},
		cli.StringFlag{Name: "tls-key-file", Value: "", Usage: "the server's private key file for TLS connection"},
	}

	for _, d := range dnames {
		flags = append(flags, storage.GetDriver(d).Flags...)
	}

	flags = append(flags, cli.BoolFlag{Name: "debug, d", Usage: "enable verbose log", EnvVar: EnvDebugName})
	app.Flags = flags
	app.Action = RunAction
	app.Run(os.Args)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
