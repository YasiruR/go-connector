package negotiation

import "net/http"

type Handler interface {
	RequestContract(w http.ResponseWriter, r *http.Request)
	OfferContract(w http.ResponseWriter, r *http.Request)
	AcceptOffer(w http.ResponseWriter, r *http.Request)
	AgreeContract(w http.ResponseWriter, r *http.Request)
	GetAgreement(w http.ResponseWriter, r *http.Request)
	VerifyAgreement(w http.ResponseWriter, r *http.Request)
	FinalizeContract(w http.ResponseWriter, r *http.Request)
	TerminateContract(w http.ResponseWriter, r *http.Request)
}
