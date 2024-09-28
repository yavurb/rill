package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) Connect(broadcastId string) (*domain.Viewer, error) {
	broadcast, err := uc.repository.GetBroadcast(broadcastId)
	if err != nil {
		return nil, domain.ErrBroadcastNotFound
	}

	viewer, err := uc.repository.CreateViewer()
	if err != nil {
		return nil, err
	}

	viewer.HandleViewer(broadcast.Track)
	broadcast.AddViewer(viewer)

	return viewer, nil
}
