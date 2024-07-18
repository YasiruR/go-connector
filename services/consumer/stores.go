package consumer

import (
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg/store/memory"
	"github.com/YasiruR/connector/protocols/negotiation"
)

// does this need a state machine?

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

type providerStore struct {
	store core.Store
}

func newProviderStore() *providerStore {
	return &providerStore{store: memory.NewStore()}
}

func (p *providerStore) add(providerId, addr string) {
	_ = p.store.Set(providerId, addr)
}

func (p *providerStore) addr(providerId string) (string, error) {
	val, err := p.store.Get(providerId)
	if err != nil {
		return ``, fmt.Errorf("get provider address failed - %v", err)
	}
	return val.(string), nil
}
