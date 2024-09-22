package signaling

import (
	"fmt"

	"github.com/pion/webrtc/v4"
)

func HandleViewer(viewerSDPChan string, track *webrtc.TrackLocalStaticRTP, viewerLocalSDPChan chan<- string) {
	broadcastTrack := track

	fmt.Println("Local track available...")

	ICEServers := []webrtc.ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
		{URLs: []string{"stun:stun1.l.google.com:19302"}},
		{URLs: []string{"stun:stun2.l.google.com:19302"}},
	}
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: ICEServers,
	}

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
