package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) Get(id string) (*domain.BroadcastSession, error) {
	broadcast, err := uc.repository.GetBroadcast(id)
	if err != nil {
		// TODO: Return an error from the domain
		return nil, err
	}

	return broadcast, nil
}
