package main

import (
	"models"

	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/otsimo/api/apipb"
)

type CatalogHeader struct {
	Title     string `toml:"title"`
	VisibleAt string `toml:"visible_at"`
	ValidDays int    `toml:"valid_days"`
}

type CatalogItem struct {
	GameID string `toml:"game_id"`
}

type CatalogFile struct {
	Catalog         CatalogHeader `toml:"catalog"`
	Featured        []CatalogItem `toml:"featured"`
	New             []CatalogItem `toml:"new"`
	Popular         []CatalogItem `toml:"popular"`
	RecentlyUpdated []CatalogItem `toml:"updated"`
}

func readCatalogFile(fpath string) (*CatalogFile, error) {
	cnf := &CatalogFile{Catalog: CatalogHeader{}}
	_, err := toml.DecodeFile(fpath, cnf)
	if err != nil {
		return nil, err
	}
	return cnf, nil
}

func toMilliseconds(s time.Time) int64 {
	return s.Unix()*1000 + int64(s.Nanosecond()/1e6)
}

func (cf *CatalogFile) Request() (*apipb.Catalog, error) {
	req := &apipb.Catalog{
		Title:     cf.Catalog.Title,
		CreatedAt: models.MillisecondsNow(),
	}
	t1, e := time.Parse(
		time.RFC3339,
		fmt.Sprintf("%s:00+00:00", cf.Catalog.VisibleAt))

	if e != nil {
		return nil, fmt.Errorf("invalid time format:%+v", e)
	}
	req.VisibleAt = toMilliseconds(t1)
	req.ExpiresAt = toMilliseconds(t1.AddDate(0, 0, cf.Catalog.ValidDays))

	for i, v := range cf.Featured {
		req.Items = append(req.Items, &apipb.CatalogItem{
			GameId:   v.GameID,
			Index:    int32(i),
			Category: apipb.CatalogCategory_FEATURED,
		})
	}
	for i, v := range cf.New {
		req.Items = append(req.Items, &apipb.CatalogItem{
			GameId:   v.GameID,
			Index:    int32(i),
			Category: apipb.CatalogCategory_NEW,
		})
	}
	for i, v := range cf.Popular {
		req.Items = append(req.Items, &apipb.CatalogItem{
			GameId:   v.GameID,
			Index:    int32(i),
			Category: apipb.CatalogCategory_POPULAR,
		})
	}
	for i, v := range cf.RecentlyUpdated {
		req.Items = append(req.Items, &apipb.CatalogItem{
			GameId:   v.GameID,
			Index:    int32(i),
			Category: apipb.CatalogCategory_RECENTLY_UPDATED,
		})
	}
	return req, nil
}
