package repository

import (
	"errors"

	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type localRepository struct {
	broadcasts []*domain.BroadcastSession
}

func NewLocalRepository(broadcasts []*domain.BroadcastSession) domain.BroadcastsRepository {
	return &localRepository{
		broadcasts,
	}
}

func (r *localRepository) GetBroadcast(id string) (*domain.BroadcastSession, error) {
	for _, broadcast := range r.broadcasts {
		if broadcast.ID == id {
			return broadcast, nil
		}
	}

	return nil, errors.New("could not get broadcast")
}

func (r *localRepository) GetBroadcasts() ([]*domain.BroadcastSession, error) {
	return r.broadcasts, nil
}
