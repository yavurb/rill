package application

import (
	"log"

	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) DeleteViewer(id string) error {
	viewer, err := uc.repository.GetViewer(id)
	if err != nil {
		return domain.ErrViewerNotFound
	}

	broadcast, err := uc.repository.GetBroadcast(viewer.BroadcastID)
	if err != nil {
		return domain.ErrBroadcastNotFound
	}

	log.Printf("Deleting viewer with ID %s\n", id)
	err = uc.repository.DeleteViewer(id)
	if err != nil {
		return err
	}

	broadcast.RemoveViewer(viewer)

	return nil
}
