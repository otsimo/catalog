package mongodb

import (
	"errors"
	"models"

	"gopkg.in/mgo.v2/bson"
	mgo "gopkg.in/mgo.v2"
)

func (d *MongoDBDriver) GetCurrent() (*models.Catalog, error) {
	return nil, errors.New("Not Found")
}

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

func (d *MongoDBDriver) List() ([]*models.Catalog, error) {
	c := d.Session.DB("").C(CatalogCollection)
	var result []*models.Catalog
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
