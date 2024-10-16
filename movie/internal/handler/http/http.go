package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"sen1or/micromovie/movie/internal/controller/movie"
)

type Handler struct {
	ctrl movie.Controller
}

func New(ctrl movie.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	movieID := r.FormValue("id")
	details, err := h.ctrl.Get(r.Context(), movieID)

	if err != nil && errors.Is(err, movie.MetadataErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
