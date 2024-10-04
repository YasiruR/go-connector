package policy

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/models/odrl"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
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
		return odrl.Offer{}, stores.QueryFailed(`policy`, `get`, err)
	}

	if val == nil {
		return odrl.Offer{}, stores.InvalidKey(id)
	}

	return val.(odrl.Offer), nil
}
