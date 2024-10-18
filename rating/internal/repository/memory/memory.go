package memory

import (
	"context"
	"sen1or/micromovie/rating/internal/repository"
	model "sen1or/micromovie/rating/pkg/model"
	"sync"
)

type Repository struct {
	mu   *sync.RWMutex
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func NewRepository() *Repository {
	return &Repository{
		data: map[model.RecordType]map[model.RecordID][]model.Rating{},
	}
}

func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	ratings, ok := r.data[recordType][recordID]

	if !ok {
		return nil, repository.ErrNotFound
	}

	return ratings, nil
}

func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)
	return nil
}
