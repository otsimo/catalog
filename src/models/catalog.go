package models

import (
	"time"

	"github.com/otsimo/api/apipb"
	"gopkg.in/mgo.v2/bson"
)

type Catalog struct {
	Id        bson.ObjectId `bson:"_id"`
	Title     string        `bson:"title"`
	Catalog   bson.Binary   `bson:"catalog"`
	CreatedAt int64         `bson:"created_at,omitempty"`
}

func NewCatalogModel(c *apipb.Catalog) (*Catalog, error) {
	b, err := c.Marshal()
	if err != nil {
		return nil, err
	}
	return &Catalog{
		Id:    bson.NewObjectId(),
		Title: c.Title,
		Catalog: bson.Binary{
			Kind: 0x80,
			Data: b,
		},
		CreatedAt: MillisecondsNow(),
	}, nil
}

func MillisecondsNow() int64 {
	s := time.Now()
	return s.Unix()*1000 + int64(s.Nanosecond()/1e6)
}
