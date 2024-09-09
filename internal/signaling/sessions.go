package signaling

import "github.com/pion/webrtc/v4"

var Broadcasts []*BroadcastSession

type BroadcastSession struct {
	Track *webrtc.TrackLocalStaticRTP
	ID    string
}
