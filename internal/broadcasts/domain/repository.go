package domain

import "github.com/pion/webrtc/v4"

type BroadcastsRepository interface {
	GetBroadcast(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	CreateBroadcast(remoteSDPSession, localSDPSession, broadcastTitle string, track <-chan *webrtc.TrackLocalStaticRTP) (*BroadcastSession, error)
}
