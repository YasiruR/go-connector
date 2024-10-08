package postgresql

import (
	"github.com/YasiruR/go-connector/domain/api/exchanger/http/postgresql"
	"github.com/YasiruR/go-connector/domain/data-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/pkg/middleware"
	"net/http"
)

type Handler struct {
	exchanger data_plane.Exchanger
}

func NewHandler(e data_plane.Exchanger) *Handler {
	return &Handler{exchanger: e}
}

func (h *Handler) HandlePush(w http.ResponseWriter, r *http.Request) {
	var req postgresql.PushRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`push transfer request`,
			err)), http.StatusBadRequest)
		return
	}

	if err := h.exchanger.PushWithCredentials(data_plane.DatabasePostgresql, req.Database); err != nil {
		middleware.WriteError(w, errors.TransferFailed(data_plane.DatabasePostgresql, `push`,
			req.Endpoint, err), http.StatusInternalServerError)
		return
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`push transfer request`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandlePull(w http.ResponseWriter, r *http.Request) {
	var req postgresql.PullRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`pull transfer request`,
			err)), http.StatusBadRequest)
		return
	}

	if err := h.exchanger.PullWithCredentials(data_plane.DatabasePostgresql, req.Database); err != nil {
		middleware.WriteError(w, errors.TransferFailed(data_plane.DatabasePostgresql, `pull`,
			req.Endpoint, err), http.StatusInternalServerError)
		return
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`pull transfer request`,
			err)), http.StatusInternalServerError)
	}
}
