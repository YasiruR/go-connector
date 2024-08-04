package catalog

import "net/http"

type Handler interface {
	HandleCatalogRequest(w http.ResponseWriter, r *http.Request)
	HandleDatasetRequest(w http.ResponseWriter, r *http.Request)
}
