package webrtc

import "github.com/pion/webrtc/v4"

var (
	Broadcasts []*BroadcasterSession
)

type BroadcasterSession struct {
	ID    string
	Track *webrtc.TrackLocalStaticRTP
}
