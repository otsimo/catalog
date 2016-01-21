package catalog

import (
	"net"
	"os"
	"storage"

	pb "github.com/otsimo/api/apipb"

	"errors"
	"models"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/mgo.v2/bson"
)

type Server struct {
	Config  *Config
	Storage storage.Driver
	Oidc    *Client
}

func (s *Server) ListenGRPC() {
	grpcPort := s.Config.GetGrpcPortString()
	//Listen
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("server.go: failed to listen %v for grpc", err)
	}
	var l = &log.Logger{
		Out:       os.Stdout,
		Formatter: &log.TextFormatter{FullTimestamp: true},
		Hooks:     make(log.LevelHooks),
		Level:     log.GetLevel(),
	}
	grpclog.SetLogger(l)

	var opts []grpc.ServerOption
	if s.Config.TlsCertFile != "" && s.Config.TlsKeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(s.Config.TlsCertFile, s.Config.TlsKeyFile)
		if err != nil {
			log.Fatalf("server.go: Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	catalogGrpc := &catalogGrpcServer{
		server: s,
	}

	pb.RegisterCatalogServiceServer(grpcServer, catalogGrpc)
	log.Infof("server.go: Binding %s for grpc", grpcPort)
	//Serve
	grpcServer.Serve(lis)
}

func NewServer(config *Config, driver storage.Driver) *Server {
	server := &Server{
		Config:  config,
		Storage: driver,
	}
	c, err := NewOIDCClient(config.ClientID, config.ClientSecret, config.AuthDiscovery)
	if err != nil {
		log.Fatal("Unable to create Oidc client", err)
	}
	server.Oidc = c
	return server
}

func (s *Server) Insert(c *pb.Catalog, email string, id bson.ObjectId) error {
	if c == nil {
		return errors.New("catalog is null")
	}
	mc, err := models.NewCatalogModel(c, email, id)
	if err != nil {
		return err
	}
	old, err := s.Storage.GetByTitle(mc.Title)
	if err != models.ErrorNotFound {
		if err == nil {
			if old.Status == pb.CatalogStatus_APPROVED {
				return fmt.Errorf("cannot update approved catalog")
			}
			old.Sync(mc)
			old.Status = pb.CatalogStatus_DRAFT
			return s.Storage.Update(old)
		}
		return err
	}
	mc.Status = pb.CatalogStatus_DRAFT
	return s.Storage.Put(mc)
}

func (s *Server) Approve(title string) error {
	return s.Storage.ChangeStatus(title, pb.CatalogStatus_APPROVED)
}

func (s *Server) Current() (*pb.Catalog, error) {
	query := pb.CatalogListRequest{
		Limit:  1,
		Status: pb.CatalogListRequest_ONLY_APPROVED,
		Time:   models.MillisecondsNow(),
	}
	res, err := s.Storage.List(query)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, models.ErrorNotFound
	}
	return res[0].ToProto(), nil
}
