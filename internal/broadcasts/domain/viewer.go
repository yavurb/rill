package domain

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/pkg/utils"
)

type ViewerCreate struct {
	BroadcastID string
}

type ViewerEvent struct {
	Response chan<- string
	Data     any
	Event    string
}

type Viewer struct {
	ctx      context.Context
	EventOut chan ViewerEvent
	EventIn  chan ViewerEvent
	cancel   context.CancelCauseFunc

	BroadcastID string
	ID          string
}

func (v *Viewer) ListenEvent() <-chan ViewerEvent {
	return v.EventOut
}

func (v *Viewer) SendEvent(event ViewerEvent) {
	v.EventIn <- event
}

func (v *Viewer) Close(cause error) {
	if v.cancel == nil {
		return
	}

	v.cancel(cause)
}

func (v *Viewer) ContextClose() <-chan struct{} {
	if v.ctx == nil {
		return nil
	}

	return v.ctx.Done()
}

func (v *Viewer) HandleViewer(track *webrtc.TrackLocalStaticRTP, config *config.Config) {
	ctx, cancel := context.WithCancelCause(context.Background())
	v.ctx = ctx
	v.cancel = cancel

	go func() {
		ICEServers := []webrtc.ICEServer{}
		for _, server := range config.WebRTC.IceServers {
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

			log.Println("Viewer - Sending candidate...")
			v.EventOut <- ViewerEvent{Event: "candidate", Data: string(outbound)}
		})

	Viewer:
		for {
			select {
			case <-ctx.Done():
				log.Println("Viewer connection closed")
				break Viewer
			case event := <-v.EventIn:
				switch event.Event {
				case "candidate":
					log.Println("Viewer - Candidate received")
					if err := peerConnection.AddICECandidate(event.Data.(webrtc.ICECandidateInit)); err != nil {
						log.Println("Error adding ICE candidate: ", err)
						cancel(err)
						break Viewer
					}
				case "offer":
					log.Println("Viewer - Offer received")
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
}
