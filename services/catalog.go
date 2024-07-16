package services

import (
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg/store/memory"
	"sync"
)

type Catalog struct {
	store  core.Store
	states *sync.Map
}

func NewCatalog() *Catalog {
	mockStore := memory.NewStore()

	return &Catalog{store: mockStore, states: new(sync.Map)}
}

func (c *Catalog) GetCatalog(filter any) {}

func (c *Catalog) GetDataset(id string) {}
