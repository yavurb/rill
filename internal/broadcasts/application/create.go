package application

import (
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/signaling"
)

func (uc *usecase) Create(broadcastTitle string) (*domain.BroadcastSession, error) {
	trackChan := make(chan *webrtc.TrackLocalStaticRTP)
	localSDPSessionChan := make(chan string)
	broadcastEventChanIn := make(chan domain.BroadcastEvent)
	broadcastEventChanOut := make(chan domain.BroadcastEvent)
	broadcastEventChanOut2 := make(chan domain.BroadcastEvent)

	ctx, cancel := signaling.HandleBroadcasterConnection(broadcastEventChanIn, broadcastEventChanOut, broadcastEventChanOut2, trackChan, localSDPSessionChan)

	broadcast, err := uc.repository.CreateBroadcast(domain.BroadcastCreate{
		Title:              broadcastTitle,
		BroadcastEventIn:   broadcastEventChanIn,
		BroadcastEventOut:  broadcastEventChanOut,
		BroadcastEventOut2: broadcastEventChanOut2,
		Ctx:                ctx,
		Cancel:             cancel,
	})
	if err != nil {
		cancel(err)
		return nil, err
	}

	go broadcast.SetTrack(trackChan)

	return broadcast, nil
}
