package grpc

import (
	"context"
	"errors"
	"sen1or/micromovie/gen"
	"sen1or/micromovie/rating/internal/controller/rating"
	"sen1or/micromovie/rating/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Error(codes.InvalidArgument, "missing information")
	}

	r, err := h.ctrl.GetAggregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{
		RatingValue: float32(r),
	}, nil
}
func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.UserId == "" || req.RecordId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing information")
	}

	rating := &model.Rating{
		UserID:     model.UserID(req.UserId),
		RecordID:   model.RecordID(req.RecordId),
		RecordType: model.RecordType(req.RecordType),
		Value:      model.RatingValue(req.RatingValue),
	}

	if err := h.ctrl.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), rating); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.PutRatingResponse{}, nil
}
