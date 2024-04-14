package domain

import "github.com/pion/webrtc/v4"

type BroadcastSession struct {
	Track            *webrtc.TrackLocalStaticRTP
	ID               string
	RemoteSDPSession string
	LocalSDPSession  string
}
