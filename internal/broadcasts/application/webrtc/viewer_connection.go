package webrtc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/broadcasts/domain"
	"github.com/yavurb/rill/internal/pkg/utils"
)

type viewerConnectionUsecase struct {
	viewer *domain.Viewer
	config *config.Config
	logger domain.Logger
}

func NewViewerConnectionUsecase(config *config.Config, viewer *domain.Viewer, logger domain.Logger) *viewerConnectionUsecase {
	return &viewerConnectionUsecase{
		config: config,
		viewer: viewer,
		logger: logger,
	}
}

func (uc *viewerConnectionUsecase) Connect(track *webrtc.TrackLocalStaticRTP) error {
	ctx, cancel := context.WithCancelCause(context.Background())
	uc.viewer.SetContext(ctx, cancel)

	go func() {
		ICEServers := []webrtc.ICEServer{}
		for _, server := range uc.config.WebRTC.IceServers {
			ICEServers = append(ICEServers, webrtc.ICEServer{
				URLs:       server.Urls,
				Username:   server.Username,
				Credential: server.Credential,
			})
		}

		peerConnectionConfig := webrtc.Configuration{
			ICEServers: ICEServers,
		}

		// Create a new PeerConnection
		peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
		if err != nil {
			panic(err)
		}

		rtpSender, err := peerConnection.AddTrack(track)
		if err != nil {
			panic(err)
		}

		// Read incoming RTCP packets
		// Before these packets are returned they are processed by interceptors. For things
		// like NACK this needs to be called.
		go func() {
			rtcpBuf := make([]byte, 1500)
			for {
				if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
					return
				}
			}
		}()

		peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			if candidate == nil {
				return
			}

			outbound, marshalErr := json.Marshal(candidate.ToJSON())
			if marshalErr != nil {
				panic(marshalErr)
			}

			uc.logger.Info("Viewer - Sending candidate...")
			uc.viewer.EventOut <- domain.ViewerEvent{Event: "candidate", Data: string(outbound)}
		})

	Viewer:
		for {
			select {
			case <-ctx.Done():
				uc.logger.Info("Viewer connection closed")
				break Viewer
			case event := <-uc.viewer.EventIn:
				switch event.Event {
				case "candidate":
					uc.logger.Info("Viewer - Candidate received")
					if err := peerConnection.AddICECandidate(event.Data.(webrtc.ICECandidateInit)); err != nil {
						uc.logger.Info("Error adding ICE candidate: ", err)
						cancel(err)
						break Viewer
					}
				case "offer":
					uc.logger.Info("Viewer - Offer received")
					offer := webrtc.SessionDescription{}
					err := utils.Decode(event.Data.(string), &offer, false)
					if err != nil {
						cancel(err)
						break Viewer
					}

					// Set the remote SessionDescription
					err = peerConnection.SetRemoteDescription(offer)
					if err != nil {
						cancel(err)
						break Viewer
					}

					// Create answer
					answer, err := peerConnection.CreateAnswer(nil)
					if err != nil {
						cancel(err)
						break Viewer
					}

					// Sets the LocalDescription, and starts our UDP listeners
					err = peerConnection.SetLocalDescription(answer)
					if err != nil {
						cancel(err)
						break Viewer
					}

					if event.Response == nil {
						cancel(errors.New("response channel is nil"))
						break Viewer
					}

					base64Answer, err := utils.Encode(answer, false)
					if err != nil {
						cancel(err)
						break Viewer
					}

					event.Response <- base64Answer

				}
			}
		}
	}()

	return nil
}
