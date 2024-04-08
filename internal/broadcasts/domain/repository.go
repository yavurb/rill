package domain

type BroadcastsRepository interface {
	GetBroadcast(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
}
