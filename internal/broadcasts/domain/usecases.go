package domain

import "github.com/pion/webrtc/v4"

// FIXME: Use cases should receive a context.Context to handle timeouts and cancellations.
type BroadcastsUsecases interface {
	Get(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	Create(broadcastTitle string) (*BroadcastSession, error)
	SaveICECandidate(id string, candidate any) error
	SaveOffer(id, sdp string) (string, error)
	Delete(id string) error

	// Viewer related use cases
	Connect(broadcastId string) (*Viewer, error)
	SaveViewerICECandidate(broadcastId string, candidate any) error
	SaveViewerOffer(broadcastId, sdp string) (string, error)
	DeleteViewer(viewerId string) error
}

type BroadcastConnectionUsecase interface {
	Connect(broadcast *BroadcastSession) error
}

type ViewerConnectionUsecase interface {
	Connect(viewer *Viewer, track *webrtc.TrackLocalStaticRTP) error
}
