package dsp

import "net/http"

type HTTPServer interface {
	HandleCatalogRequest(w http.ResponseWriter, r *http.Request)
	GetNegotiation(w http.ResponseWriter, r *http.Request)
	HandleContractRequest(w http.ResponseWriter, r *http.Request)
	HandleContractAgreement(w http.ResponseWriter, r *http.Request)
}
