package application

import (
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/signaling"
)

func (uc *usecase) Create(broadcastTitle string) (*domain.BroadcastSession, error) {
	trackChan := make(chan *webrtc.TrackLocalStaticRTP)
	localSDPSessionChan := make(chan string)
	broadcastEventChan := make(chan domain.BroadcastEvent)

	ctx, cancel := signaling.HandleBroadcasterConnection(broadcastEventChan, trackChan, localSDPSessionChan)

	broadcast, err := uc.repository.CreateBroadcast(domain.BroadcastCreate{
		Title:          broadcastTitle,
		BroadcastEvent: broadcastEventChan,
		Ctx:            ctx,
		Cancel:         cancel,
	})
	if err != nil {
		cancel(err)
		return nil, err
	}

	go broadcast.SetTrack(trackChan)

	return broadcast, nil
}
