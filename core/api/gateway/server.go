package gateway

import "net/http"

type HTTPServer interface {
	CreateAsset(w http.ResponseWriter, r *http.Request)
	RequestContract(w http.ResponseWriter, r *http.Request)
	AgreeContract(w http.ResponseWriter, r *http.Request)
}
