package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/x264"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"

	_ "github.com/pion/mediadevices/pkg/driver/videotest"
)

func main() {
	x264Params, err := x264.NewParams()
	if err != nil {
		panic(err)
	}
	x264Params.BitRate = 500_000

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&x264Params),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	pc, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
	})
	if err != nil {
		panic(err)
	}

	pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s", connectionState.String())
	})

	s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatI420)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		},
		Codec: codecSelector,
	})
	if err != nil {
		panic(err)
	}

	for _, track := range s.GetTracks() {
		log.Printf("Track (ID: %s, StreamID: %s) added", track.ID(), track.StreamID())
		track.OnEnded(func(err error) {
			log.Printf("Track (ID: %s) ended with error: %v\n",
				track.ID(), err)
		})

		if _, err := pc.AddTransceiverFromTrack(track, webrtc.RtpTransceiverInit{Direction: webrtc.RTPTransceiverDirectionSendonly}); err != nil {
			panic(err)
		}
	}

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err := pc.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	resp, err := http.Post(os.Args[1], "application/sdp", strings.NewReader(offer.SDP))
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  string(body),
	}); err != nil {
		panic(err)
	}

	select {}
}
