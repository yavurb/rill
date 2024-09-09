package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/signaling"
)

func (uc *usecase) Connect(remoteSdp, broadcastId string) (*domain.Viewer, error) {
	broadcast, err := uc.repository.GetBroadcast(broadcastId)
	if err != nil {
		return nil, domain.ErrBroadcastNotFound
	}

	viewerLocalSDPChan := make(chan string)

	go signaling.HandleViewer(remoteSdp, broadcast.Track, viewerLocalSDPChan)

	viewerLocalSDP := <-viewerLocalSDPChan

	viewer := &domain.Viewer{
		Events:          make(chan string, 1),
		LocalSDPSession: viewerLocalSDP,
	}

	broadcast.AddViewer(viewer)

	return viewer, nil
}
