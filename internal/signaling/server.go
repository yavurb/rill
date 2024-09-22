package signaling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type (
	subscriber struct {
		event chan string
	}
	publisher struct {
		subscribers map[*subscriber]struct{}
	}
	serverCtx struct {
		publishers map[string]*publisher
	}
)

func NewSignalingServer(e *echo.Echo) *serverCtx {
	return &serverCtx{}
}

func HandleBroadcasterConnection(
	eventChan chan domain.BroadcastEvent,
	trackChan chan<- *webrtc.TrackLocalStaticRTP,
	broadcasterLocalSDPChan chan<- string,
) (context.Context, context.CancelCauseFunc) {
	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		defer cancel(nil)

		ICEServers := []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
			{URLs: []string{"stun:stun1.l.google.com:19302"}},
			{URLs: []string{"stun:stun2.l.google.com:19302"}},
			{URLs: []string{"stun:stun3.l.google.com:19302"}},
			{URLs: []string{"stun:stun4.l.google.com:19302"}},
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

			eventChan <- domain.BroadcastEvent{Data: "candidate", Event: string(outbound)}
		})

	Broadcast:
		for {
			select {
			case <-ctx.Done():
				log.Println("Broadcaster connection closed")
				break Broadcast
			case event := <-eventChan:
				log.Println("Event received: ", event)
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
					Decode(event.Data.(string), &offer)

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

					outbound, marshalErr := json.Marshal(answer)
					if marshalErr != nil {
						log.Println("Error marshaling answer: ", marshalErr)
						cancel(marshalErr)
						break Broadcast
					}

					eventChan <- domain.BroadcastEvent{Data: "answer", Event: Encode(outbound)}
				}
			}
		}
	}()

	return ctx, cancel
}

func HandleViewer(viewerSDPChan string, track *webrtc.TrackLocalStaticRTP, viewerLocalSDPChan chan<- string) {
	broadcastTrack := track

	fmt.Println("Local track available...")

	ICEServers := []webrtc.ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
		{URLs: []string{"stun:stun1.l.google.com:19302"}},
		{URLs: []string{"stun:stun2.l.google.com:19302"}},
		{URLs: []string{"stun:stun3.l.google.com:19302"}},
		{URLs: []string{"stun:stun4.l.google.com:19302"}},
	}
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: ICEServers,
	}

	fmt.Printf("I'm passign through here")

	recvOnlyOffer := webrtc.SessionDescription{}
	Decode(viewerSDPChan, &recvOnlyOffer)

	// Create a new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		panic(err)
	}

	rtpSender, err := peerConnection.AddTrack(broadcastTrack)
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

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(recvOnlyOffer)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one ng message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Get the LocalDescription and take it to base64 so we can paste in browser
	viewerLocalSDPChan <- fmt.Sprint(Encode(*peerConnection.LocalDescription()))

	done := make(chan bool)
	<-done
}
