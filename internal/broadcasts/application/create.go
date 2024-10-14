package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) Create(broadcastTitle string) (*domain.BroadcastSession, error) {
	broadcast, err := uc.repository.CreateBroadcast(domain.BroadcastCreate{
		Title: broadcastTitle,
	})
	if err != nil {
		uc.logger.Errorf("failed to create broadcast: %v", err)
		return nil, err
	}

	err = uc.broadcastUsecase.Connect(broadcast)
	if err != nil {
		uc.logger.Errorf("failed to make connection: %v", err)
		return nil, err
	}

	return broadcast, nil
}
