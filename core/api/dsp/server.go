package dsp

import "net/http"

type HTTPServer interface {
	GetNegotiation(w http.ResponseWriter, r *http.Request)
	HandleContractRequest(w http.ResponseWriter, r *http.Request)
}
