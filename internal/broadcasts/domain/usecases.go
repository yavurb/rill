package domain

type BroadcastsUsecases interface {
	Get(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
}
