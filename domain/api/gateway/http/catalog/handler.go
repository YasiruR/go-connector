package catalog

import "net/http"

type Handler interface {
	CreateOffer(w http.ResponseWriter, r *http.Request)
	CreateDataset(w http.ResponseWriter, r *http.Request)
	RequestCatalog(w http.ResponseWriter, r *http.Request)
	RequestDataset(w http.ResponseWriter, r *http.Request)
	GetStoredCatalogs(w http.ResponseWriter, r *http.Request)
}
