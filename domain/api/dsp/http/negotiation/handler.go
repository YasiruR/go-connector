package negotiation

import "net/http"

type Handler interface {
	HandleContractRequest(w http.ResponseWriter, r *http.Request)
	HandleContractOffer(w http.ResponseWriter, r *http.Request)
	HandleContractAgreement(w http.ResponseWriter, r *http.Request)
	HandleAgreementVerification(w http.ResponseWriter, r *http.Request)
	HandleNegotiationEvent(w http.ResponseWriter, r *http.Request)
	GetNegotiation(w http.ResponseWriter, r *http.Request)
}
