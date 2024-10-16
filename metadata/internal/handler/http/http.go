package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sen1or/micromovie/metadata/internal/controller/metadata"
	"sen1or/micromovie/metadata/internal/repository"
)

type Handler struct {
	ctrl *metadata.Controller
}

func NewHandler(ctrl *metadata.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	m, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("error while encoding metadata: %s\n", err)
	}
}
