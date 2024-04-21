package domain

type BroadcastsUsecases interface {
	Get(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	Create(remoteSDPSession, broadcastTitle string) (string, error)
	Connect(remoteSDPSession, broadcastId string) (string, error)
}
