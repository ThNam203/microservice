package metadata

import (
	"context"
	"errors"
	repository "sen1or/micromovie/metadata/internal/repository"
	model "sen1or/micromovie/metadata/pkg"
)

var (
	ErrNotFound = errors.New("metadata not found for movie")
)

type metadataRepository interface {
	Get(ctx context.Context, movieID string) (*model.Metadata, error)
	Put(ctx context.Context, movieID string, metadata *model.Metadata) error
}

type Controller struct {
	Repo metadataRepository
}

func NewController(repo metadataRepository) *Controller {
	return &Controller{
		Repo: repo,
	}
}

func (r *Controller) Get(ctx context.Context, movieID string) (*model.Metadata, error) {
	record, err := r.Repo.Get(ctx, movieID)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}

	return record, err
}
