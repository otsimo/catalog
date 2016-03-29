package models

import (
	"time"

	apipb "github.com/otsimo/otsimopb"
	"gopkg.in/mgo.v2/bson"
)

type Catalog struct {
	Id          bson.ObjectId       `bson:"_id"`
	Title       string              `bson:"title"`
	CreatedAt   int64               `bson:"created_at"`
	UpdatedAt   int64               `bson:"updated_at"`
	VisibleAt   int64               `bson:"visible_at"`
	ExpiresAt   int64               `bson:"expires_at"`
	Items       []apipb.CatalogItem `bson:"items"`
	AuthorEmail string              `bson:"author_email"`
	AuthorID    bson.ObjectId       `bson:"author_id"`
	Status      apipb.CatalogStatus `bson:"status"`
}

func NewCatalogModel(c *apipb.Catalog, email string, user_id bson.ObjectId) (*Catalog, error) {
	mc := &Catalog{
		Id:          bson.NewObjectId(),
		Title:       c.Title,
		VisibleAt:   c.VisibleAt,
		ExpiresAt:   c.ExpiresAt,
		CreatedAt:   MillisecondsNow(),
		UpdatedAt:   MillisecondsNow(),
		AuthorEmail: email,
		AuthorID:    user_id,
		Status:      c.Status,
	}
	mc.Items = make([]apipb.CatalogItem, len(c.Items))
	for i, p := range c.Items {
		mc.Items[i] = *p
	}
	return mc, nil
}

func MillisecondsNow() int64 {
	s := time.Now()
	return s.Unix()*1000 + int64(s.Nanosecond()/1e6)
}

func (mc *Catalog) Sync(other *Catalog) {
	mc.UpdatedAt = MillisecondsNow()
	mc.ExpiresAt = other.ExpiresAt
	mc.VisibleAt = other.VisibleAt
	mc.Items = make([]apipb.CatalogItem, len(other.Items))
	for i, p := range other.Items {
		mc.Items[i] = p
	}
}

func (mc *Catalog) ToProto() *apipb.Catalog {
	c := &apipb.Catalog{
		Title:     mc.Title,
		ExpiresAt: mc.ExpiresAt,
		VisibleAt: mc.VisibleAt,
		CreatedAt: mc.UpdatedAt,
		Status:    mc.Status,
	}
	c.Items = make([]*apipb.CatalogItem, len(mc.Items))
	for i, p := range mc.Items {
		c.Items[i] = &apipb.CatalogItem{
			GameId:   p.GameId,
			Category: p.Category,
			Index:    p.Index,
		}
	}
	return c
}
