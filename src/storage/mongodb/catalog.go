package mongodb

import (
	"models"

	apipb "github.com/otsimo/otsimopb"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (d *MongoDBDriver) GetById(id bson.ObjectId) (*models.Catalog, error) {
	c := d.Session.DB("").C(CatalogCollection)
	var doc models.Catalog
	err := c.FindId(id).One(&doc)
	if err == mgo.ErrNotFound {
		return nil, models.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (d *MongoDBDriver) GetByTitle(title string) (*models.Catalog, error) {
	c := d.Session.DB("").C(CatalogCollection)
	var doc models.Catalog
	err := c.Find(bson.M{"title": title}).One(&doc)
	if err == mgo.ErrNotFound {
		return nil, models.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (d *MongoDBDriver) Put(mc *models.Catalog) error {
	c := d.Session.DB("").C(CatalogCollection)
	return c.Insert(mc)
}

func (d *MongoDBDriver) Update(mc *models.Catalog) error {
	c := d.Session.DB("").C(CatalogCollection)
	return c.UpdateId(mc.Id, bson.M{"$set": mc})
}

func (d *MongoDBDriver) ChangeStatus(title string, stat apipb.CatalogStatus) error {
	c := d.Session.DB("").C(CatalogCollection)
	return c.Update(bson.M{"title": title}, bson.M{"$set": bson.M{"status": stat}})
}

func (d *MongoDBDriver) List(q apipb.CatalogListRequest) ([]*models.Catalog, error) {
	c := d.Session.DB("").C(CatalogCollection)
	var result []*models.Catalog
	query := bson.M{}
	if q.Limit == 0 {
		q.Limit = 100
	}
	switch q.Status {
	case apipb.CatalogListRequest_ONLY_DRAFT:
		query["status"] = apipb.CatalogStatus_DRAFT
	case apipb.CatalogListRequest_ONLY_APPROVED:
		query["status"] = apipb.CatalogStatus_APPROVED
	}
	if q.Time > 0 {
		query["visible_at"] = bson.M{"$lte": q.Time}
		if q.HideExpired {
			m := int64(0)
			now := models.MillisecondsNow()
			if q.Time > now {
				m = q.Time
			} else {
				m = now
			}
			query["expires_at"] = bson.M{"$gt": m}
		} else {
			query["expires_at"] = bson.M{"$gt": q.Time}
		}
	} else if q.HideExpired {
		query["expires_at"] = bson.M{"$gt": models.MillisecondsNow()}
	}

	iter := c.Find(query).Limit(int(q.Limit)).Sort("-visible_at").Iter()
	err := iter.All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
