package catalog

import (
	"github.com/YasiruR/connector/core/pkg"
)

type Catalog struct {
	store  pkg.DataStore
	states pkg.DataStore
}

func NewCatalog() *Catalog {
	return &Catalog{}
}

func (c *Catalog) GetCatalog(filter any) error {
	return nil
}

func (c *Catalog) GetDataset(id string) {}
