package catalog

import (
	"errors"

	"github.com/otsimo/api/apipb"
	"golang.org/x/net/context"
)

type catalogGrpcServer struct {
	server *Server
}

func (w *catalogGrpcServer) GetCatalog(ctx context.Context, in *apipb.CatalogRequest) (*apipb.CatalogResponse, error) {

	return nil, errors.New("not implemented")
}
