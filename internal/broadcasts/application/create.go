package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) Create(broadcastTitle string) (*domain.BroadcastSession, error) {
	broadcast, err := uc.repository.CreateBroadcast(domain.BroadcastCreate{
		Title: broadcastTitle,
	})
	if err != nil {
		return nil, err
	}

	broadcast.MakeRTCConnection(uc.config)

	return broadcast, nil
}
