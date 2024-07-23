package stores

import (
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

type Policy struct {
	store pkg.Collection
}

func NewPolicyStore(plugins core.Plugins) *Policy {
	plugins.Log.Info("initialized policy store")
	return &Policy{store: plugins.Database.NewDataStore()}
}

func (p *Policy) SetOffer(id string, val odrl.Offer) {
	_ = p.store.Set(id, val)
}

func (p *Policy) Offer(id string) (odrl.Offer, error) {
	val, err := p.store.Get(id)
	if err != nil {
		return odrl.Offer{}, errors.QueryFailed(`policy`, `get`, err)
	}
	return val.(odrl.Offer), nil
}
