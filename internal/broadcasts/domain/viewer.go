package domain

import (
	"context"
)

type ViewerCreate struct {
	BroadcastID string
}

type ViewerEvent struct {
	Response chan<- string
	Data     any
	Event    string
}

type Viewer struct {
	ctx      context.Context
	EventOut chan ViewerEvent
	EventIn  chan ViewerEvent
	cancel   context.CancelCauseFunc

	BroadcastID string
	ID          string
}

func (v *Viewer) ListenEvent() <-chan ViewerEvent {
	return v.EventOut
}

func (v *Viewer) SendEvent(event ViewerEvent) {
	v.EventIn <- event
}

func (v *Viewer) Close(cause error) {
	if v.cancel == nil {
		return
	}

	v.cancel(cause)
}

func (v *Viewer) SetContext(ctx context.Context, cancel context.CancelCauseFunc) {
	v.ctx = ctx
	v.cancel = cancel
}

func (v *Viewer) ContextClose() <-chan struct{} {
	if v.ctx == nil {
		return nil
	}

	return v.ctx.Done()
}
