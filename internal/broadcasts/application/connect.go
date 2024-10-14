package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) Connect(broadcastId string) (*domain.Viewer, error) {
	broadcast, err := uc.repository.GetBroadcast(broadcastId)
	if err != nil {
		return nil, domain.ErrBroadcastNotFound
	}

	viewer, err := uc.repository.CreateViewer(domain.ViewerCreate{BroadcastID: broadcastId})
	if err != nil {
		return nil, err
	}

	err = uc.viewerUsecase.Connect(viewer, broadcast.Track)
	if err != nil {
		uc.logger.Errorf("failed to connect viewer: %v", err)
		return nil, err
	}

	broadcast.AddViewer(viewer)

	return viewer, nil
}
