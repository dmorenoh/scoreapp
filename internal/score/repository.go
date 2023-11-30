package score

import (
	"context"
	"errors"
	"sync"
)

var ErrVersionMismatch = errors.New("version mismatch")

type Repository interface {
	Get(ctx context.Context, user uint) (*Score, error)
	Save(ctx context.Context, score *Score) error
	FindAll(ctx context.Context, filter interface{}) ([]*Score, error)
}

type inMemoryRepository struct {
	sync.RWMutex
	scores map[uint]*Entity
}

func NewInMemoryRepository(scores map[uint]*Entity) Repository {
	return &inMemoryRepository{
		scores: scores,
	}
}

func (i *inMemoryRepository) Save(ctx context.Context, score *Score) error {
	i.Lock()
	defer i.Unlock()

	existing, ok := i.scores[score.user]
	if !ok {
		i.scores[score.user] = &Entity{User: score.user, Total: score.total, Version: score.version + 1}
	} else if existing.Version != score.version {
		return ErrVersionMismatch
	} else {
		i.scores[score.user] = &Entity{User: score.user, Total: score.total, Version: score.version + 1}
	}

	return nil
}

func (i *inMemoryRepository) FindAll(ctx context.Context, filter interface{}) ([]*Score, error) {
	i.RLock()
	defer i.RUnlock()

	filteredList, err := entitiesMap(i.scores).sortValues().filter(filter)
	if err != nil {
		return nil, err
	}

	results := make([]*Score, len(filteredList))
	for i2, entity := range filteredList {
		results[i2] = &Score{
			user:    entity.User,
			total:   entity.Total,
			version: entity.Version,
		}
	}
	return results, nil
}

func (i *inMemoryRepository) Get(ctx context.Context, user uint) (*Score, error) {
	i.RLock()
	defer i.RUnlock()

	entity := i.scores[user]
	if entity == nil {
		return nil, nil
	}
	return &Score{
		user:    entity.User,
		total:   entity.Total,
		version: entity.Version,
	}, nil
}
