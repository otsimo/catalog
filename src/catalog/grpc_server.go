package catalog

import (
	"errors"
	"models"

	"github.com/Sirupsen/logrus"
	"github.com/otsimo/api/apipb"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

type catalogGrpcServer struct {
	server *Server
}

func (w *catalogGrpcServer) Pull(ctx context.Context, in *apipb.CatalogPullRequest) (*apipb.Catalog, error) {
	return w.server.Current()
}

func (w *catalogGrpcServer) Push(ctx context.Context, in *apipb.Catalog) (*apipb.Response, error) {
	jwt, err := getJWTToken(ctx)
	if err != nil {
		logrus.Errorf("grpc_server.go: failed to get jwt %+v", err)
		return nil, errors.New("failed to get jwt")
	}
	id, email, err := w.authToken(jwt, true)
	if err != nil {
		logrus.Errorf("grpc_server.go: failed to authorize user %+v", err)
		return nil, errors.New("unauthorized user")
	}
	if !bson.IsObjectIdHex(id) {
		return nil, models.ErrorInvalidID
	}

	err = w.server.Insert(in, email, bson.ObjectIdHex(id))
	if err != nil {
		return nil, err
	}
	return &apipb.Response{Type: 0, Message: "success"}, nil
}

func (w *catalogGrpcServer) Approve(ctx context.Context, in *apipb.CatalogApproveRequest) (*apipb.Response, error) {
	jwt, err := getJWTToken(ctx)
	if err != nil {
		logrus.Errorf("grpc_server.go: failed to get jwt %+v", err)
		return nil, errors.New("failed to get jwt")
	}
	_, _, err = w.authToken(jwt, true)
	if err != nil {
		logrus.Errorf("grpc_server.go: failed to authorize user %+v", err)
		return nil, errors.New("unauthorized user")
	}
	err = w.server.Approve(in.Title)
	if err != nil {
		return nil, err
	}
	return &apipb.Response{Type: 0, Message: "success"}, nil
}

func (w *catalogGrpcServer) List(ctx context.Context, query *apipb.CatalogListRequest) (*apipb.CatalogListResponse, error) {
	res, err := w.server.Storage.List(*query)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, models.ErrorNotFound
	}
	result := make([]*apipb.Catalog, len(res))
	for i, p := range res {
		result[i] = p.ToProto()
	}
	return &apipb.CatalogListResponse{
		Catalogs: result,
	}, nil
}
