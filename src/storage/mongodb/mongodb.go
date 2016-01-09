package mongodb

import (
	"storage"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	mgo "gopkg.in/mgo.v2"
)

const (
	MongoDBDriverName string = "mongodb"
	mongoURLFlag      string = "mongodb-url"

	CatalogCollection string = "Catalog"
)

func init() {
	storage.Register(MongoDBDriverName, &storage.RegisteredDriver{
		New: newMongoDriver,
		Flags: []cli.Flag{
			cli.StringFlag{Name: mongoURLFlag, Value: "mongodb://localhost:27017/Otsimo", Usage: "MongoDB url", EnvVar: "MONGODB_URL"},
		},
	})
}

func newMongoDriver(ctx *cli.Context) (storage.Driver, error) {
	url := ctx.String(mongoURLFlag)

	s, err := mgo.Dial(url)

	if err != nil {
		return nil, err
	}
	log.Info("mongodb.go: connected to mongodb")
	md := &MongoDBDriver{
		Session: s,
	}
	return md, nil
}

type MongoDBDriver struct {
	Session *mgo.Session
}

func (d MongoDBDriver) Name() string {
	return MongoDBDriverName
}

func millisecondsNow() int64 {
	s := time.Now()
	return s.Unix()*1000 + int64(s.Nanosecond()/1e6)
}
