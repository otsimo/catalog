package main

import (
	"catalog"
	"fmt"
	"os"
	"storage"
	_ "storage/mongodb"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var Version string
var config *catalog.Config = catalog.NewConfig()

func RunAction(c *cli.Context) error {
	config.Debug = c.Bool("debug")
	config.GrpcPort = c.Int("grpc-port")
	config.HealthPort = c.Int("health-port")
	config.TlsCertFile = c.String("tls-cert-file")
	config.TlsKeyFile = c.String("tls-key-file")
	config.ClientID = c.String("client-id")
	config.ClientSecret = c.String("client-secret")
	config.AuthDiscovery = c.String("discovery")

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}
	//find driver name
	sname := c.String("storage")
	if sname == "" || sname == "none" {
		cli.ShowAppHelp(c)
		return fmt.Errorf("main.go: storage flag is missing or it cannot be 'none'")
	}

	//get driver
	driver := storage.GetDriver(sname)
	if driver == nil {
		return fmt.Errorf("main.go: storage driver '%s' not found\n", sname)
	}

	//load storage driver
	s, err := driver.New(c)
	if err != nil {
		return fmt.Errorf("main.go: error while creating new storage driver[%s]:%v", s, err)
	}

	server := catalog.NewServer(config, s)

	ch := server.Oidc.SyncProviderConfig(config.AuthDiscovery)
	defer func() {
		// stop the background process
		ch <- struct{}{}
	}()
	return server.ListenGRPC()
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
	app.Name = "otsimo-catalog"
	app.Version = Version
	app.Usage = "Otsimo Catalog Service"
	app.Author = "Sercan Degirmenci <sercan@otsimo.com>"
	dnames := storage.GetDriverNames()
	var flags []cli.Flag

	flags = []cli.Flag{
		cli.IntFlag{Name: "grpc-port", Value: config.GrpcPort, Usage: "grpc server port"},
		cli.IntFlag{Name: "health-port", Value: config.HealthPort, Usage: "health check port"},
		cli.StringFlag{Name: "storage, s", Value: "none", Usage: fmt.Sprintf("the storage driver. Available drivers: %s", strings.Join(dnames, ", "))},
		cli.StringFlag{Name: "tls-cert-file", Value: "", Usage: "the server's certificate file for TLS connection"},
		cli.StringFlag{Name: "tls-key-file", Value: "", Usage: "the server's private key file for TLS connection"},
		cli.StringFlag{Name: "client-id", Value: "", Usage: "client id"},
		cli.StringFlag{Name: "client-secret", Value: "", Usage: "client secret"},
		cli.StringFlag{Name: "discovery", Value: "https://connect.otsimo.com", Usage: "auth discovery url"},
	}
	flags = withEnvs("OTSIMO_CATALOG", flags)
	for _, d := range dnames {
		flags = append(flags, storage.GetDriver(d).Flags...)
	}

	flags = append(flags, cli.BoolFlag{Name: "debug, d", Usage: "enable verbose log", EnvVar: "OTSIMO_CATALOG_DEBUG"})
	app.Flags = flags
	app.Action = RunAction

	log.Infoln("running", app.Name, "version:", app.Version)
	app.Run(os.Args)
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
