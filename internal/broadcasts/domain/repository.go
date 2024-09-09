package domain

type BroadcastsRepository interface {
	GetBroadcast(id string) (*BroadcastSession, error)
	GetBroadcasts() ([]*BroadcastSession, error)
	CreateBroadcast(broadcast BroadcastCreate) (*BroadcastSession, error)
	UpdateBroadcast(id string, broadcast BroadcastUpdate) (*BroadcastSession, error)
	DeleteBroadcast(id string) error
}
