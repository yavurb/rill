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

	broadcast.SendEvent(domain.BroadcastEvent{
		Event: "offer",
		Data:  sdp,
	})

	var answer string
	for broadcastEvent := range broadcast.ListenEvent2() {
		log.Printf("SaveOffer: %s\n", broadcastEvent.Event)
		if broadcastEvent.Event != "answer" {
			continue
		}

		answer = broadcastEvent.Data.(string)
		break
	}

	log.Println("Final answer", answer)
	return answer, nil
}
