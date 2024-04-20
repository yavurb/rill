package application

import (
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/webrtc"
)

func (uc *usecase) Connect(remoteSdp, broadcastId string) (string, error) {
	broadcast, err := uc.repository.GetBroadcast(broadcastId)
	if err != nil {
		return "", domain.ErrNotFound
	}

	viewerLocalSDPChan := make(chan string)

	go webrtc.HandleViewer(remoteSdp, broadcast.Track, viewerLocalSDPChan)

	viewerLocalSDP := <-viewerLocalSDPChan

	return viewerLocalSDP, nil
}
