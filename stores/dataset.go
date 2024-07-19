package stores

import "github.com/YasiruR/connector/core/models/dcat"

type Dataset struct{}

func (d *Dataset) Set(id string, val dcat.Dataset) {}

func (d *Dataset) Get(id string) (dcat.Dataset, error) {
	return dcat.Dataset{}, nil
}
