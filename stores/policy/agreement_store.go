package policy

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/models/odrl"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

const (
	collAgreement            = `agreement`
	collNegotiationAgreement = `negotiation-agreement`
)

type Agreement struct {
	agrColl  pkg.Collection
	cnAgrMap pkg.Collection
}

func NewAgreementStore(plugins domain.Plugins) *Agreement {
	plugins.Log.Info("initialized agreement store")
	return &Agreement{agrColl: plugins.Database.NewCollection(), cnAgrMap: plugins.Database.NewCollection()}
}

func (a *Agreement) AddAgreement(cnId string, val odrl.Agreement) {
	_ = a.agrColl.Set(val.Id, val)
	a.setNegotiationId(cnId, val.Id)
}

func (a *Agreement) setNegotiationId(cnId, agrId string) {
	_ = a.cnAgrMap.Set(cnId, agrId)
}

func (a *Agreement) Agreement(id string) (odrl.Agreement, error) {
	val, err := a.agrColl.Get(id)
	if err != nil {
		return odrl.Agreement{}, stores.QueryFailed(collAgreement, `Get`, err)
	}

	if val == nil {
		return odrl.Agreement{}, stores.InvalidKey(id)
	}

	return val.(odrl.Agreement), nil
}

//func (a *Agreement) AgreementByNegotiationID(cnId string) (odrl.Agreement, error) {
//	agrId, err := a.cnAgrMap.Get(cnId)
//	if err != nil {
//		return odrl.Agreement{}, errors.QueryFailed(collNegotiationAgreement, `Get`, err)
//	}
//
//	if agrId == nil {
//		return odrl.Agreement{}, errors.InvalidKey(cnId)
//	}
//
//	return a.Agreement(agrId.(string))
//}
