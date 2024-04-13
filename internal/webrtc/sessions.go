package webrtc

import "github.com/pion/webrtc/v4"

var (
	Broadcasts []*BroadcastSession
)

type BroadcastSession struct {
	ID    string
	Track *webrtc.TrackLocalStaticRTP
}
