package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sen1or/micromovie/rating/internal/controller/rating"
	model "sen1or/micromovie/rating/pkg"
	"strconv"
)

type Handler struct {
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	recordID := model.RecordID(r.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(r.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		sum, err := h.ctrl.GetAggregatedRating(r.Context(), recordID, recordType)
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(sum); err != nil {
			log.Printf("response encode error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPut:
		userId := model.UserID(r.FormValue("userId"))
		value, err := strconv.ParseFloat(r.FormValue("value"), 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h.ctrl.PutRating(r.Context(), recordID, recordType, &model.Rating{RecordID: recordID, RecordType: recordType, UserID: userId, Value: model.RatingValue(value)}); err != nil {
			log.Printf("failed to put rating: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
