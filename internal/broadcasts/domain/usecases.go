package domain

type BroadcastsUsecases interface {
	Get(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	Create(remoteSDPSession string) (string, error)
	Connect(remoteSDPSession, broadcastId string) (string, error)
}
