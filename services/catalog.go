package services

import (
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg/store/memory"
)

type Catalog struct {
	store  core.Store
	states core.Store
}

func NewCatalog() *Catalog {
	return &Catalog{store: memory.NewStore(), states: memory.NewStore()}
}

func (c *Catalog) GetCatalog(filter any) error {
	return nil
}

func (c *Catalog) GetDataset(id string) {}
