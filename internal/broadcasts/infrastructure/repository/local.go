package repository

import (
	"errors"

	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/pkg/publicid"
)

type localRepository struct {
	broadcasts []*domain.BroadcastSession
}

const broadcastIdPrefix = "br"

func NewLocalRepository(broadcasts []*domain.BroadcastSession) domain.BroadcastsRepository {
	return &localRepository{
		broadcasts,
	}
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

func (r *localRepository) CreateBroadcast(remoteSDPSession, localSDPSession, broadcastTitle string, TrackChan <-chan *webrtc.TrackLocalStaticRTP) (*domain.BroadcastSession, error) {
	track := <-TrackChan
	broadcastId, err := publicid.New(broadcastIdPrefix, 12)
	if err != nil {
		return nil, err
	}

	broadcast := &domain.BroadcastSession{
		ID:    broadcastId,
		Title: broadcastTitle,
		Track: track, RemoteSDPSession: remoteSDPSession,
		LocalSDPSession: localSDPSession,
	}

	r.broadcasts = append(r.broadcasts, broadcast)

	return broadcast, nil
}
