package storage

import (
	"fmt"

	"models"

	"github.com/codegangsta/cli"
	apipb "github.com/otsimo/otsimopb"

	"gopkg.in/mgo.v2/bson"
)

type Driver interface {
	Name() string
	GetById(bson.ObjectId) (*models.Catalog, error)
	GetByTitle(string) (*models.Catalog, error)
	List(apipb.CatalogListRequest) ([]*models.Catalog, error)
	Put(*models.Catalog) error
	Update(*models.Catalog) error
	ChangeStatus(string, apipb.CatalogStatus) error
}

type RegisteredDriver struct {
	New   func(*cli.Context) (Driver, error)
	Flags []cli.Flag
}

var drivers map[string]*RegisteredDriver

func init() {
	drivers = make(map[string]*RegisteredDriver)
}

func Register(name string, rd *RegisteredDriver) error {
	if _, ext := drivers[name]; ext {
		return fmt.Errorf("Name already registered %s", name)
	}
	drivers[name] = rd
	return nil
}

func GetDriverNames() []string {
	drives := make([]string, 0)

	for name, _ := range drivers {
		drives = append(drives, name)
	}
	return drives
}

func GetDriver(name string) *RegisteredDriver {
	return drivers[name]
}
