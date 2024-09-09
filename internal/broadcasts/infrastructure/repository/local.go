package repository

import (
	"errors"
	"sync"

	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/pkg/publicid"
)

type localRepository struct {
	broadcasts      []*domain.BroadcastSession
	broadcastsMutex sync.Mutex
}

const broadcastIdPrefix = "br"

func NewLocalRepository() domain.BroadcastsRepository {
	return &localRepository{}
}

func (r *localRepository) GetBroadcast(id string) (*domain.BroadcastSession, error) {
	for _, broadcast := range r.broadcasts {
		if broadcast.ID == id {
			return broadcast, nil
		}
	}

	return nil, errors.New("could not get broadcast")
}

func (r *localRepository) GetBroadcasts() ([]*domain.BroadcastSession, error) {
	return r.broadcasts, nil
}

func (r *localRepository) CreateBroadcast(broadcast domain.BroadcastCreate) (*domain.BroadcastSession, error) {
	broadcastId, err := publicid.New(broadcastIdPrefix, 12)
	if err != nil {
		return nil, err
	}

	broadcast_ := &domain.BroadcastSession{
		ID:               broadcastId,
		Title:            broadcast.Title,
		LocalSDPSession:  broadcast.LocalSDPSession,
		RemoteSDPSession: broadcast.RemoteSDPSession,
		Viewers:          make(map[*domain.Viewer]struct{}),
	}

	broadcast_.SetCtx(broadcast.Ctx, broadcast.Cancel)

	r.broadcastsMutex.Lock()
	r.broadcasts = append(r.broadcasts, broadcast_)
	r.broadcastsMutex.Unlock()

	return broadcast_, nil
}

func (r *localRepository) UpdateBroadcast(id string, broadcast domain.BroadcastUpdate) (*domain.BroadcastSession, error) {
	broadcast_, err := r.GetBroadcast(id)
	if err != nil {
		return nil, err
	}

	track := <-broadcast.Track

	broadcast_.Track = track
	broadcast_.Title = broadcast.Title
	broadcast_.LocalSDPSession = broadcast.LocalSDPSession
	broadcast_.RemoteSDPSession = broadcast.RemoteSDPSession

	return broadcast_, nil
}
