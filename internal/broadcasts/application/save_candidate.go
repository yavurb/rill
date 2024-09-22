package application

import (
	"encoding/json"

	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) SaveICECandidate(id string, candidate any) error {
	broadcast, err := uc.repository.GetBroadcast(id)
	if err != nil {
		return domain.ErrBroadcastNotFound
	}

	var candidate_ webrtc.ICECandidateInit

	jsonCandidate, _ := json.Marshal(candidate)
	if err := json.Unmarshal(jsonCandidate, &candidate_); err != nil {
		return err
	}

	broadcast.SendEvent(domain.BroadcastEvent{
		Event: "candidate",
		Data:  candidate_,
	})

	return nil
}
