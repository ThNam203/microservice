package grpc

import (
	"context"
	"errors"
	"sen1or/micromovie/gen"
	"sen1or/micromovie/metadata/pkg/model"
	"sen1or/micromovie/movie/internal/controller/movie"
	"sen1or/micromovie/movie/internal/grpcutil"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gen.UnimplementedMovieServiceServer
	ctrl movie.Controller
}

func New(ctrl movie.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing information")
	}

	res, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, movie.MetadataErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Metadata: model.MetadataToProto(&res.Metadata),
			Rating:   float32(*res.Rating),
		},
	}, nil
}
