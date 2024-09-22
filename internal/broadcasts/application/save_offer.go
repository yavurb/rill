package application

import (
	"log"

	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) SaveOffer(id, sdp string) (string, error) {
	broadcast, err := uc.repository.GetBroadcast(id)
	if err != nil {
		return "", domain.ErrBroadcastNotFound
	}

	response := make(chan string)
	broadcast.SendEvent(domain.BroadcastEvent{
		Event:    "offer",
		Data:     sdp,
		Response: response,
	})

	answer := <-response
	log.Println("Final answer", answer)
	return answer, nil
}
