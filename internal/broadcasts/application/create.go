package application

import (
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/signaling"
)

func (uc *usecase) Create(remoteSDPSession, broadcastTitle string) (*domain.BroadcastSession, error) {
	trackChan := make(chan *webrtc.TrackLocalStaticRTP)
	localSDPSessionChan := make(chan string)

	ctx, cancel := signaling.HandleBroadcasterConnection(remoteSDPSession, trackChan, localSDPSessionChan)

	broadcastLocalSDPSession := <-localSDPSessionChan

	broadcast, err := uc.repository.CreateBroadcast(domain.BroadcastCreate{
		Title:            broadcastTitle,
		RemoteSDPSession: remoteSDPSession,
		LocalSDPSession:  broadcastLocalSDPSession,
		Ctx:              ctx,
		Cancel:           cancel,
	})
	if err != nil {
		cancel(err)
		return nil, err
	}

	go uc.repository.UpdateBroadcast(broadcast.ID, domain.BroadcastUpdate{
		Title:            broadcast.Title,
		RemoteSDPSession: broadcast.RemoteSDPSession,
		LocalSDPSession:  broadcast.LocalSDPSession,
		Track:            trackChan,
	})

	return broadcast, nil
}
