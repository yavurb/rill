package application

import (
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

	broadcastEvent := <-broadcast.ListenEvent()

	if broadcastEvent.Event != "answer" {
		return "", domain.ErrBroadcastInvalidEvent
	}

	return broadcastEvent.Data.(string), nil
}
