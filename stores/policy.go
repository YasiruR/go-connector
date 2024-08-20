package stores

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
)

// Policy is a store that exists within a Provider to persist any created policy
type Policy struct {
	store pkg.Collection
}

func NewPolicyStore(plugins domain.Plugins) *Policy {
	plugins.Log.Info("initialized policy store")
	return &Policy{store: plugins.Database.NewCollection()}
}

func (p *Policy) AddOffer(id string, val odrl.Offer) {
	_ = p.store.Set(id, val)
}

func (p *Policy) Offer(id string) (odrl.Offer, error) {
	val, err := p.store.Get(id)
	if err != nil {
		return odrl.Offer{}, errors.QueryFailed(`policy`, `get`, err)
	}
	return val.(odrl.Offer), nil
}
