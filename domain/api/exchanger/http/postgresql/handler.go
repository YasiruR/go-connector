package postgresql

import "net/http"

type Handler interface {
	HandlePush(w http.ResponseWriter, r *http.Request)
	HandlePull(w http.ResponseWriter, r *http.Request)
}
