package transfer

import (
	"net/http"
)

type Handler interface {
	HandleTransfers(w http.ResponseWriter, r *http.Request)
	HandleTransferRequest(w http.ResponseWriter, r *http.Request)
	HandleTransferStart(w http.ResponseWriter, r *http.Request)
	HandleTransferSuspension(w http.ResponseWriter, r *http.Request)
	HandleTransferCompletion(w http.ResponseWriter, r *http.Request)
}
