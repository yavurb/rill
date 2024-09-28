package application

import (
	"encoding/json"

	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

func (uc *usecase) SaveViewerICECandidate(id string, candidate any) error {
	viewer, err := uc.repository.GetViewer(id)
	if err != nil {
		return domain.ErrViewerNotFound
	}

	var candidate_ webrtc.ICECandidateInit

	jsonCandidate, _ := json.Marshal(candidate)
	if err := json.Unmarshal(jsonCandidate, &candidate_); err != nil {
		return err
	}

	viewer.SendEvent(domain.ViewerEvent{
		Event: "candidate",
		Data:  candidate_,
	})

	return nil
}
