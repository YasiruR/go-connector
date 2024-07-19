package gateway

import "net/http"

type HTTPServer interface {
	CreateAsset(w http.ResponseWriter, r *http.Request)
	InitContractRequest(w http.ResponseWriter, r *http.Request)
	HandleContractAgreement(w http.ResponseWriter, r *http.Request)
}
