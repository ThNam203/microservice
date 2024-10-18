package repository

import (
	"context"
	model "sen1or/micromovie/metadata/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func NewRepository() *Repository {
	return &Repository{
		data: map[string]*model.Metadata{},
	}
}

func (r *Repository) Get(_ context.Context, movieID string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()

	x, ok := r.data[movieID]
	if !ok {
		return nil, ErrNotFound
	}

	return x, nil
}

func (r *Repository) Put(_ context.Context, movieID string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()

	r.data[movieID] = metadata
	return nil
}
