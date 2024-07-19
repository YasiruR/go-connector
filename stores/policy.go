package stores

import (
	"fmt"
	"github.com/YasiruR/connector/core/models/odrl"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/pkg/store/memory"
)

type Policy struct {
	store pkg.Store
}

func NewPolicyStore() *Policy {
	return &Policy{
		store: memory.NewStore(), // todo inject store as a dependency
	}
}

func (p *Policy) SetOffer(id string, val odrl.Offer) {
	_ = p.store.Set(id, val)
}

func (p *Policy) GetOffer(id string) (odrl.Offer, error) {
	val, err := p.store.Get(id)
	if err != nil {
		return odrl.Offer{}, fmt.Errorf("store failure - %w", err)
	}
	return val.(odrl.Offer), nil
}
