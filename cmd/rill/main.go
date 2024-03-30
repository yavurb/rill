package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pion/webrtc/v4"
	lwebrtc "github.com/yavurb/rill/internal/webrtc"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4321"},
	}))
	e.Use(middleware.Logger())

	localTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	broadcasterSDPChan := make(chan string)
	broadcasterLocalSDPChan := make(chan string)
	viewerLocalSDPChan := make(chan string)

	go lwebrtc.HandleBroadcasterConnection(broadcasterSDPChan, broadcasterLocalSDPChan)
	// go test(viewerSDPChan, viewerLocalSDPChan)

	e.POST("/broadcaster", func(c echo.Context) error {
		type sdpS struct {
			SDP string `json:"sdp"`
		}
		type response struct {
			SDP string `json:"sdp"`
		}

		var s sdpS
		err := c.Bind(&s)
		if err != nil {
			return err
		}

		broadcasterSDPChan <- s.SDP

		sdp := <-broadcasterLocalSDPChan

		return c.JSON(http.StatusOK, response{SDP: sdp})
	})

	e.POST("/viewer", func(c echo.Context) error {
		type sdpS struct {
			SDP string `json:"sdp"`
		}
		type response struct {
			SDP string `json:"sdp"`
		}

		var s sdpS
		err := c.Bind(&s)
		if err != nil {
			return err
		}

		fmt.Println("Init endpoint")
		fmt.Printf("Trach Channel has %d data stored", len(localTrackChan))

		if lwebrtc.CurrentSession.Track == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "No broadcast available",
			})
		}

		go lwebrtc.HandleViewer(s.SDP, lwebrtc.CurrentSession.Track, viewerLocalSDPChan)

		select {
		case sdp := <-viewerLocalSDPChan:
			return c.JSON(http.StatusOK, response{SDP: sdp})
		case <-time.After(time.Second * 10):
			return c.JSON(http.StatusRequestTimeout, echo.Map{
				"error": "timeout waiting for SDP",
			})
		}

	})

	e.Logger.Fatal(e.Start(":8910"))
}
