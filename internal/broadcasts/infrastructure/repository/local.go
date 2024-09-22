package repository

import (
	"errors"
	"sync"

	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/pkg/publicid"
)

type localRepository struct {
	broadcasts      map[string]*domain.BroadcastSession
	broadcastsMutex sync.Mutex
}

const broadcastIdPrefix = "br"

func NewLocalRepository() domain.BroadcastsRepository {
	return &localRepository{
		broadcasts: make(map[string]*domain.BroadcastSession),
	}
}

func (r *localRepository) GetBroadcast(id string) (*domain.BroadcastSession, error) {
	broadcast, ok := r.broadcasts[id]
	if !ok {
		return nil, errors.New("could not get broadcast")
	}

	return broadcast, nil
}

func (r *localRepository) GetBroadcasts() ([]*domain.BroadcastSession, error) {
	broadcasts := make([]*domain.BroadcastSession, 0, len(r.broadcasts))

	for _, broadcast := range r.broadcasts {
		broadcasts = append(broadcasts, broadcast)
	}

	return broadcasts, nil
}

func (r *localRepository) CreateBroadcast(broadcast domain.BroadcastCreate) (*domain.BroadcastSession, error) {
	broadcastID, err := publicid.New(broadcastIdPrefix, 12)
	if err != nil {
		return nil, err
	}

	broadcast_ := &domain.BroadcastSession{
		ID:      broadcastID,
		Title:   broadcast.Title,
		Event:   broadcast.BroadcastEvent,
		Viewers: make(map[*domain.Viewer]struct{}),
	}

	broadcast_.SetCtx(broadcast.Ctx, broadcast.Cancel)

	r.broadcastsMutex.Lock()
	r.broadcasts[broadcastID] = broadcast_
	r.broadcastsMutex.Unlock()

	return broadcast_, nil
}

func (r *localRepository) UpdateBroadcast(id string, broadcast domain.BroadcastUpdate) (*domain.BroadcastSession, error) {
	broadcast_, err := r.GetBroadcast(id)
	if err != nil {
		return nil, err
	}

	if broadcast.Title != "" {
		broadcast_.Title = broadcast.Title
	}

	return broadcast_, nil
}

func (r *localRepository) DeleteBroadcast(id string) error {
	r.broadcastsMutex.Lock()
	delete(r.broadcasts, id)
	r.broadcastsMutex.Unlock()

	return nil
}
