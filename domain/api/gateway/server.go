package gateway

import "net/http"

type HTTPServer interface {
	CreatePolicy(w http.ResponseWriter, r *http.Request)
	CreateDataset(w http.ResponseWriter, r *http.Request)
	RequestCatalog(w http.ResponseWriter, r *http.Request)
	RequestDataset(w http.ResponseWriter, r *http.Request)
	RequestContract(w http.ResponseWriter, r *http.Request)
	AgreeContract(w http.ResponseWriter, r *http.Request)
	VerifyAgreement(w http.ResponseWriter, r *http.Request)
	FinalizeContract(w http.ResponseWriter, r *http.Request)
	GetAgreement(w http.ResponseWriter, r *http.Request)
	Start()
}
