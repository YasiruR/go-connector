package provider

import (
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg/store/memory"
	"github.com/YasiruR/connector/protocols/negotiation"
)

type stateStore struct {
	store core.Store
}

func newStateStore() *stateStore {
	return &stateStore{store: memory.NewStore()}
}

func (s *stateStore) get(id string) (negotiation.State, error) {
	val, err := s.store.Get(id)
	if err != nil {
		return ``, fmt.Errorf("get state failed - %v", err)
	}
	return val.(negotiation.State), nil
}

func (s *stateStore) set(id string, val negotiation.State) {
	_ = s.store.Set(id, val)
}
