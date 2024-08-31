package policy

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
)

const (
	agreementCollection            = `agreement-collection`
	agreementNegotiationCollection = `agreement-negotiation-collection`
)

type Agreement struct {
	store pkg.Collection
	cnMap pkg.Collection
}

func NewAgreementStore(plugins domain.Plugins) *Agreement {
	plugins.Log.Info("initialized agreement store")
	return &Agreement{store: plugins.Database.NewCollection(), cnMap: plugins.Database.NewCollection()}
}

func (a *Agreement) AddAgreement(cnId string, val odrl.Agreement) {
	_ = a.store.Set(val.Id, val)
	a.setNegotiationId(cnId, val.Id)
}

func (a *Agreement) setNegotiationId(cnId, agrId string) {
	_ = a.cnMap.Set(cnId, agrId)
}

func (a *Agreement) Agreement(id string) (odrl.Agreement, error) {
	val, err := a.store.Get(id)
	if err != nil {
		return odrl.Agreement{}, errors.QueryFailed(agreementCollection, `Get`, err)
	}
	return val.(odrl.Agreement), nil
}

func (a *Agreement) AgreementByNegotiationID(cnId string) (odrl.Agreement, error) {
	agrId, err := a.cnMap.Get(cnId)
	if err != nil {
		return odrl.Agreement{}, errors.QueryFailed(agreementNegotiationCollection, `Get`, err)
	}

	return a.Agreement(agrId.(string))
}
