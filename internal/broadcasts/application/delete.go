package application

import (
	"log"

	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) Delete(id string) error {
	_, err := uc.repository.GetBroadcast(id)
	if err != nil {
		return domain.ErrBroadcastNotFound
	}

	log.Printf("Deleting broadcast with ID %s\n", id)
	uc.repository.DeleteBroadcast(id)

	return nil
}
