package stores

import (
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

type Policy struct {
	db pkg.Database
}

func NewPolicyStore(db pkg.Database) *Policy {
	return &Policy{db: db}
}

func (p *Policy) SetOffer(id string, val odrl.Offer) {
	_ = p.db.Set(id, val)
}

func (p *Policy) GetOffer(id string) (odrl.Offer, error) {
	val, err := p.db.Get(id)
	if err != nil {
		return odrl.Offer{}, errors.QueryFailed(`get`, err)
	}
	return val.(odrl.Offer), nil
}
