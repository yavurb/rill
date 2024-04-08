package application

import "github.com/yavurb/rill/internal/broadcasts/domain"

func (uc *usecase) GetBroadcasts() ([]*domain.BroadcastSession, error) {
	broadcasts, err := uc.repository.GetBroadcasts()
	if err != nil {
		// TODO: Return a domain error
		return nil, err
	}

	return broadcasts, nil
}
