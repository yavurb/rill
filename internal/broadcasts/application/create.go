package application

import (
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/signaling"
)

func (uc *usecase) Create(remoteSDPSession, broadcastTitle string) (string, error) {
	trackChan := make(chan *webrtc.TrackLocalStaticRTP)
	localSDPSessionChan := make(chan string)

	go signaling.HandleBroadcasterConnection(remoteSDPSession, trackChan, localSDPSessionChan)

	broadcastLocalSDPSession := <-localSDPSessionChan

	go uc.repository.CreateBroadcast(remoteSDPSession, broadcastLocalSDPSession, broadcastTitle, trackChan)

	return broadcastLocalSDPSession, nil
}
