package domain

import "github.com/pion/webrtc/v4"

type BroadcastsRepository interface {
	GetBroadcast(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	CreateBroadcast(remoteSDPSession, localSDPSession string, track *webrtc.TrackLocalStaticRTP) (*BroadcastSession, error)
}
