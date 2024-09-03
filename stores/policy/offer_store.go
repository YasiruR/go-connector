package policy

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
)

// OfferStore is a store that exists within a Provider to persist any created policy
type OfferStore struct {
	store pkg.Collection
}

func NewOfferStore(plugins domain.Plugins) *OfferStore {
	plugins.Log.Info("initialized offer store")
	return &OfferStore{store: plugins.Database.NewCollection()}
}

func (o *OfferStore) AddOffer(id string, val odrl.Offer) {
	_ = o.store.Set(id, val)
}

func (o *OfferStore) Offer(id string) (odrl.Offer, error) {
	val, err := o.store.Get(id)
	if err != nil {
		return odrl.Offer{}, core.QueryFailed(`policy`, `get`, err)
	}

	if val == nil {
		return odrl.Offer{}, core.InvalidKey(id)
	}

	return val.(odrl.Offer), nil
}
