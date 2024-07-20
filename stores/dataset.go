package stores

import (
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/dcat"
)

type Dataset struct {
	store pkg.Database
}

func NewDatasetStore(store pkg.Database) *Dataset {
	return &Dataset{store: store}
}

func (d *Dataset) Set(id string, val dcat.Dataset) {
	_ = d.store.Set(id, val)
}

func (d *Dataset) Get(id string) (dcat.Dataset, error) {
	val, err := d.store.Get(id)
	if err != nil {
		return dcat.Dataset{}, errors.QueryFailed(`get`, err)
	}
	return val.(dcat.Dataset), nil
}
