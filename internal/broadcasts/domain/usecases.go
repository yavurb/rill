package domain

// FIXME: Use cases should receive a context.Context to handle timeouts and cancellations.
type BroadcastsUsecases interface {
	Get(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	Create(broadcastTitle string) (*BroadcastSession, error)
	SaveICECandidate(id string, candidate any) error
	SaveOffer(id, sdp string) (string, error)
	Connect(remoteSDPSession, broadcastId string) (*Viewer, error)
	Delete(id string) error
}
