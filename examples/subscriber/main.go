package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/pion/webrtc/v3"

	_ "github.com/pion/mediadevices/pkg/driver/audiotest"
	_ "github.com/pion/mediadevices/pkg/driver/videotest"
)

func main() {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
	})
	if err != nil {
		panic(err)
	}

	pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s", connectionState.String())
	})

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Track (ID: %s, StreamID: %s) added", track.ID(), track.StreamID())
		for {
			p, _, err := track.ReadRTP()
			if err != nil {
				if err == io.EOF {
					return
				}
				panic(err)
			}
			log.Printf("Got RTP packet for track %s of size %d", track.ID(), len(p.Payload))
		}
	})

	resp, err := http.Post(os.Args[1], "application/sdp", bytes.NewReader([]byte{}))
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  string(body),
	}); err != nil {
		panic(err)
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)

	if err := pc.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	<-gatherComplete

	location, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		panic(err)
	}

	if _, err := http.DefaultClient.Do(&http.Request{
		Method: http.MethodPatch,
		URL:    location,
		Header: http.Header{
			"Content-Type": []string{"application/sdp"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(pc.LocalDescription().SDP))),
	}); err != nil {
		panic(err)
	}

	select {}
}
