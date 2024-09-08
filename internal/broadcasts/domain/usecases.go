package domain

type BroadcastsUsecases interface {
	Get(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	Create(remoteSDPSession, broadcastTitle string) (*BroadcastSession, error)
	Connect(remoteSDPSession, broadcastId string) (*Viewer, error)
}
