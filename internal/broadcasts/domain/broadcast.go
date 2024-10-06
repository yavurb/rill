package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/config"
	"github.com/yavurb/rill/internal/signaling"
)

type BroadcastCreate struct {
	Title string
}

type BroadcastUpdate struct {
	Title string
}

type BroadcastEvent struct {
	Response chan<- string
	Data     any
	Event    string
}

type BroadcastSession struct {
	ctx      context.Context
	Track    *webrtc.TrackLocalStaticRTP
	Viewers  map[*Viewer]struct{}
	EventOut chan BroadcastEvent
	EventIn  chan BroadcastEvent
	cancel   context.CancelCauseFunc

	ID    string
	Title string

	viewersMutex sync.Mutex
}

func (b *BroadcastSession) ListenEvent() <-chan BroadcastEvent {
	return b.EventOut
}

func (b *BroadcastSession) SendEvent(event BroadcastEvent) {
	b.EventIn <- event
}

func (b *BroadcastSession) Close(cause error) {
	if b.cancel == nil {
		return
	}

	b.cancel(cause)
}

func (b *BroadcastSession) ContextClose() <-chan struct{} {
	if b.ctx == nil {
		return nil
	}

	return b.ctx.Done()
}

func (b *BroadcastSession) AddViewer(viewer *Viewer) {
	b.viewersMutex.Lock()
	b.Viewers[viewer] = struct{}{}
	b.viewersMutex.Unlock()
}

func (b *BroadcastSession) RemoveViewer(viewer *Viewer) {
	b.viewersMutex.Lock()
	delete(b.Viewers, viewer)
	b.viewersMutex.Unlock()
}

func (b *BroadcastSession) SendEventToViewers(event ViewerEvent) {
	b.viewersMutex.Lock()
	for viewer := range b.Viewers {
		viewer.EventIn <- event
	}
	b.viewersMutex.Unlock()
}

func (b *BroadcastSession) SetTrack(trackChan <-chan *webrtc.TrackLocalStaticRTP) {
	track := <-trackChan
	b.Track = track
}

func (b *BroadcastSession) MakeRTCConnection(config *config.Config) {
	ctx, cancel := context.WithCancelCause(context.Background())
	b.ctx = ctx
	b.cancel = cancel

	trackChan := make(chan *webrtc.TrackLocalStaticRTP)
	go b.SetTrack(trackChan)

	go func() {
		defer cancel(nil)

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

		m := &webrtc.MediaEngine{}
		if err := m.RegisterDefaultCodecs(); err != nil {
			// NOTE: Should we return a custom error?
			cancel(err)
			return
		}

		// Create a InterceptorRegistry. This is the user configurable RTP/RTCP Pipeline.
		// This provides NACKs, RTCP Reports and other features. If you use `webrtc.NewPeerConnection`
		// this is enabled by default. If you are manually managing You MUST create a InterceptorRegistry
		// for each PeerConnection.
		i := &interceptor.Registry{}

		// Use the default set of Interceptors
		if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
			cancel(err)
			return
		}

		// Register a intervalpli factory
		// This interceptor sends a PLI every 3 seconds. A PLI causes a video keyframe to be generated by the sender.
		// This makes our video seekable and more error resilent, but at a cost of lower picture quality and higher bitrates
		// A real world application should process incoming RTCP packets from viewers and forward them to senders
		intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
		if err != nil {
			cancel(err)
			return
		}
		i.Add(intervalPliFactory)

		// Create a new RTCPeerConnection
		peerConnection, err := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i)).NewPeerConnection(peerConnectionConfig)
		if err != nil {
			cancel(err)
			return
		}
		defer func() {
			log.Println("Closing peer connection...")
			if cErr := peerConnection.Close(); cErr != nil {
				fmt.Printf("cannot close peerConnection: %v\n", cErr)
			}
		}()

		// Allow us to receive 1 video track
		if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
			cancel(err)
			return
		}

		// Set a handler for when a new remote track starts, this just distributes all our packets
		// to connected peers
		peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
			// Create a local track, all our SFU clients will be fed via this track
			localTrack, newTrackErr := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "video", "pion")
			if newTrackErr != nil {
				panic(newTrackErr)
			}
			trackChan <- localTrack

			rtpBuf := make([]byte, 1400)
			for {
				i, _, readErr := remoteTrack.Read(rtpBuf)
				if readErr != nil {
					cancel(err)
					return
				}

				// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
				if _, err = localTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
					cancel(err)
					return
				}
			}
		})

		peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			if candidate == nil {
				return
			}

			outbound, marshalErr := json.Marshal(candidate.ToJSON())
			if marshalErr != nil {
				panic(marshalErr)
			}

			log.Println("Sending candidate...")
			b.EventOut <- BroadcastEvent{Event: "candidate", Data: string(outbound)}
		})

	Broadcast:
		for {
			select {
			case <-ctx.Done():
				log.Println("Broadcaster connection closed")
				break Broadcast
			case event := <-b.EventIn:
				switch event.Event {
				case "candidate":
					log.Println("Candidate received")
					if err := peerConnection.AddICECandidate(event.Data.(webrtc.ICECandidateInit)); err != nil {
						log.Println("Error adding ICE candidate: ", err)
						cancel(err)
						break Broadcast
					}
				case "offer":
					log.Println("Offer received")
					offer := webrtc.SessionDescription{}
					signaling.Decode(event.Data.(string), &offer)

					// Set the remote SessionDescription
					err = peerConnection.SetRemoteDescription(offer)
					if err != nil {
						cancel(err)
						break Broadcast
					}

					// Create answer
					answer, err := peerConnection.CreateAnswer(nil)
					if err != nil {
						cancel(err)
						break Broadcast
					}

					// Sets the LocalDescription, and starts our UDP listeners
					err = peerConnection.SetLocalDescription(answer)
					if err != nil {
						cancel(err)
						break Broadcast
					}

					if event.Response == nil {
						cancel(errors.New("response channel is nil"))
						break Broadcast
					}

					event.Response <- signaling.Encode(answer)
				}
			}
		}
	}()
}
