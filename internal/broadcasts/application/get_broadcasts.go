package application

import "github.com/yavurb/rill/internal/broadcasts/domain"

func (uc *usecase) GetBroadcasts() ([]*domain.BroadcastSession, error) {
	broadcasts, err := uc.repository.GetBroadcasts()
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return broadcasts, nil
}
