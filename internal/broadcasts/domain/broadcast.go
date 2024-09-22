package domain

import (
	"context"
	"sync"

	"github.com/pion/webrtc/v4"
)

type BroadcastEvent struct {
	Response chan<- string
	Data     any
	Event    string
}
type BroadcastSession struct {
	ctx      context.Context
	Track    *webrtc.TrackLocalStaticRTP
	Viewers  map[*Viewer]struct{}
	EventOut <-chan BroadcastEvent
	EventIn  chan<- BroadcastEvent
	cancel   context.CancelCauseFunc

	ID    string
	Title string

	viewersMutex sync.Mutex
}

// NOTE: Should this have a method for closing the broadcast session?
type Viewer struct {
	Events          chan<- string
	LocalSDPSession string
}

func (b *BroadcastSession) SetCtx(ctx context.Context, cancel context.CancelCauseFunc) {
	b.ctx = ctx
	b.cancel = cancel
}

func (b *BroadcastSession) SetTrack(trackChan <-chan *webrtc.TrackLocalStaticRTP) {
	track := <-trackChan
	b.Track = track
}

func (b *BroadcastSession) AddIceCandidate(candidate webrtc.ICECandidateInit) {
}

func (b *BroadcastSession) ListenEvent() <-chan BroadcastEvent {
	return b.EventOut
}

func (b *BroadcastSession) SendEvent(event BroadcastEvent) {
	b.EventIn <- event
}

func (b *BroadcastSession) Close(cause error) {
	if b.cancel == nil {
		return
	}

	b.cancel(cause)
}

func (b *BroadcastSession) ContextClose() <-chan struct{} {
	if b.ctx == nil {
		return nil
	}

	return b.ctx.Done()
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
	BroadcastEventIn  chan BroadcastEvent
	BroadcastEventOut chan BroadcastEvent
	Ctx               context.Context
	Cancel            context.CancelCauseFunc
	Title             string
}

type BroadcastUpdate struct {
	Title string
}
