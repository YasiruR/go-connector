package transfer

import "net/http"

type Handler interface {
	HandleTransferRequest(w http.ResponseWriter, r *http.Request)
}
