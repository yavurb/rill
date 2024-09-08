package domain

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

type BroadcastSession struct {
	Track            *webrtc.TrackLocalStaticRTP
	Viewers          map[*Viewer]struct{}
	ID               string
	Title            string
	RemoteSDPSession string
	LocalSDPSession  string
	viewersMutex     sync.Mutex
}

// NOTE: Should this have a method for closing the broadcast session?
type Viewer struct {
	Events          chan<- string
	LocalSDPSession string
}

func (b *BroadcastSession) AddViewer(viewer *Viewer) {
	b.viewersMutex.Lock()
	b.Viewers[viewer] = struct{}{}
	b.viewersMutex.Unlock()
}

func (b *BroadcastSession) RemoveViewer(viewer *Viewer) {
	b.viewersMutex.Lock()
	delete(b.Viewers, viewer)
	b.viewersMutex.Unlock()
}

func (b *BroadcastSession) BroadcastEvent(event string) {
	b.viewersMutex.Lock()
	for viewer := range b.Viewers {
		viewer.Events <- event
	}
	b.viewersMutex.Unlock()
}

type BroadcastCreate struct {
	Title            string
	RemoteSDPSession string
	LocalSDPSession  string
}

type BroadcastUpdate struct {
	Track            <-chan *webrtc.TrackLocalStaticRTP
	Title            string
	RemoteSDPSession string
	LocalSDPSession  string
}
