package catalog

import (
	"errors"
	"models"

	"github.com/Sirupsen/logrus"
	apipb "github.com/otsimo/otsimopb"
	"golang.org/x/net/context"
)

type catalogGrpcServer struct {
	server *Server
}

func (w *catalogGrpcServer) Pull(ctx context.Context, in *apipb.CatalogPullRequest) (*apipb.Catalog, error) {
	logrus.Debugf("[Pull]: %v", in)
	return w.server.Current()
}

func (w *catalogGrpcServer) Push(ctx context.Context, in *apipb.Catalog) (*apipb.Response, error) {
	logrus.Debugf("[Push]: %v", in)
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

	err = w.server.Insert(in, email, id)
	if err != nil {
		return nil, err
	}
	return &apipb.Response{Type: 0, Message: "success"}, nil
}

func (w *catalogGrpcServer) Approve(ctx context.Context, in *apipb.CatalogApproveRequest) (*apipb.Response, error) {
	logrus.Debugf("[Approve]: %v", in)
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
	logrus.Debugf("[List]: %v", query)
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
