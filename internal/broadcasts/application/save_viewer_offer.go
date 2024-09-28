package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) SaveViewerOffer(id, sdp string) (string, error) {
	viewer, err := uc.repository.GetViewer(id)
	if err != nil {
		return "", domain.ErrBroadcastNotFound
	}

	response := make(chan string)
	viewer.SendEvent(domain.ViewerEvent{
		Event:    "offer",
		Data:     sdp,
		Response: response,
	})

	answer := <-response
	return answer, nil
}
