package stores

import "github.com/YasiruR/connector/domain/models/odrl"

/*
	Data associated with ODRL policies are stored here.
*/

// OfferStore maintains offers published by a Provider with corresponding rules
// (e.g. permissions, prohibitions, duties) associated with a target, i.e. a dataset.
type OfferStore interface {
	AddOffer(id string, val odrl.Offer)
	Offer(id string) (odrl.Offer, error)
}

// AgreementStore stores agreements resulted by successfully concluded contract
// negotiation processes between a provider and a consumer.
type AgreementStore interface {
	// AddAgreement stores contract agreement with agreement ID as the key
	AddAgreement(id string, val odrl.Agreement)
	// Agreement retrieves contract agreement by agreement ID
	Agreement(id string) (odrl.Agreement, error)
	//AgreementByNegotiationID(cnId string) (odrl.Agreement, error)
}
