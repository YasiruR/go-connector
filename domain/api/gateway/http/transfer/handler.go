package transfer

import "net/http"

type Handler interface {
	RequestTransfer(w http.ResponseWriter, r *http.Request)
	StartTransfer(w http.ResponseWriter, r *http.Request)
	SuspendTransfer(w http.ResponseWriter, r *http.Request)
	CompleteTransfer(w http.ResponseWriter, r *http.Request)
}
