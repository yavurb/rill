package domain

import "github.com/pion/webrtc/v4"

type BroadcastSession struct {
	ID    string
	Track *webrtc.TrackLocalStaticRTP
}
