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

	webRTCConnection := NewWebRTCConnectionUsecase(uc.config, broadcast)
	err = webRTCConnection.MakeConnection()
	if err != nil {
		return nil, err
	}

	return broadcast, nil
}
