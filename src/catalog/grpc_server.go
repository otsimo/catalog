package catalog

import (
	"errors"

	"github.com/otsimo/api/apipb"
	"golang.org/x/net/context"
)

type catalogGrpcServer struct {
	server *Server
}

func (w *catalogGrpcServer) Pull(context.Context, *apipb.CatalogPullRequest) (*apipb.Catalog, error) {

	return nil, errors.New("not implemented")
}

func (w *catalogGrpcServer) Push(context.Context, *apipb.Catalog) (*apipb.Response, error) {

	return nil, errors.New("not implemented")
}
