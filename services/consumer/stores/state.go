package stores

import (
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg/store/memory"
	"github.com/YasiruR/connector/protocols/negotiation"
)

type States struct {
	store core.Store
}

func NewStateStore() *States {
	return &States{store: memory.NewStore()}
}

func (s States) Get(id string) (negotiation.State, error) {
	val, err := s.store.Get(id)
	if err != nil {
		return ``, fmt.Errorf("get state failed - %v", err)
	}
	return val.(negotiation.State), nil
}

func (s States) Set(id string, val negotiation.State) {
	_ = s.store.Set(id, val)
}
