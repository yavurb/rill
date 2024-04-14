package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) Get(id string) (*domain.BroadcastSession, error) {
	broadcast, err := uc.repository.GetBroadcast(id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return broadcast, nil
}
