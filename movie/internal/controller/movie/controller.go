package movie

import (
	"context"
	"errors"
	metadatamodel "sen1or/micromovie/metadata/pkg/model"
	"sen1or/micromovie/movie/internal/gateway"
	model "sen1or/micromovie/movie/pkg"
	ratingmodel "sen1or/micromovie/rating/pkg/model"
)

var (
	MetadataErrNotFound = errors.New("metadata not found for movie")
)

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, movieID string) (*metadatamodel.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{
		ratingGateway:   ratingGateway,
		metadataGateway: metadataGateway,
	}
}

func (c *Controller) Get(ctx context.Context, movieID string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, movieID)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, MetadataErrNotFound
	} else if err != nil {
		return nil, err
	}

	result := &model.MovieDetails{
		Metadata: *metadata,
	}

	aggregatedRating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(movieID), ratingmodel.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		// proceeds even if there is no rating
	} else if err != nil {
		return nil, err
	} else {
		result.Rating = &aggregatedRating
	}

	return result, nil
}
