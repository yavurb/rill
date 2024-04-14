package application

import (
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
	lwebrtc "github.com/yavurb/rill/internal/webrtc"
)

func (uc *usecase) Create(remoteSDPSession string) (*domain.BroadcastSession, error) {
	trackChan := make(chan *webrtc.TrackLocalStaticRTP)
	localSDPSessionChan := make(chan string)

	lwebrtc.HandleBroadcasterConnection(remoteSDPSession, trackChan, localSDPSessionChan)

	broascastLocalSDPSession := <-localSDPSessionChan
	broadcastTrack := <-trackChan

	broadcast, err := uc.repository.CreateBroadcast(remoteSDPSession, broascastLocalSDPSession, broadcastTrack)
	if err != nil {
		return nil, domain.ErrUnknown
	}

	return broadcast, nil
}
