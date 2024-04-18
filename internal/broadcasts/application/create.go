package application

import (
	"github.com/pion/webrtc/v4"
	lwebrtc "github.com/yavurb/rill/internal/webrtc"
)

func (uc *usecase) Create(remoteSDPSession string) (string, error) {
	trackChan := make(chan *webrtc.TrackLocalStaticRTP)
	localSDPSessionChan := make(chan string)

	go lwebrtc.HandleBroadcasterConnection(remoteSDPSession, trackChan, localSDPSessionChan)

	broadcastLocalSDPSession := <-localSDPSessionChan

	go uc.repository.CreateBroadcast(remoteSDPSession, broadcastLocalSDPSession, trackChan)

	return broadcastLocalSDPSession, nil
}
