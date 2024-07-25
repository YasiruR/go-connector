package dsp

import "net/http"

type CatalogHandler interface {
	HandleCatalogRequest(w http.ResponseWriter, r *http.Request)
	HandleDatasetRequest(w http.ResponseWriter, r *http.Request)
}

type NegotiationHandler interface {
	HandleContractRequest(w http.ResponseWriter, r *http.Request)
	HandleContractAgreement(w http.ResponseWriter, r *http.Request)
	HandleAgreementVerification(w http.ResponseWriter, r *http.Request)
	HandleEventConsumer(w http.ResponseWriter, r *http.Request)
	GetNegotiation(w http.ResponseWriter, r *http.Request)
}
