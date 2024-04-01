package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pion/webrtc/v4"
	"github.com/yavurb/rill/internal/pkg/publicid"
	lwebrtc "github.com/yavurb/rill/internal/webrtc"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4321"},
	}))
	e.Use(middleware.Logger())

	localTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	broadcasterLocalSDPChan := make(chan string)
	viewerLocalSDPChan := make(chan string)

	// go test(viewerSDPChan, viewerLocalSDPChan)

	e.POST("/broadcaster", func(c echo.Context) error {
		type sdpS struct {
			SDP string `json:"sdp"`
		}
		type response struct {
			SDP string `json:"sdp"`
		}

		ranId, err := publicid.New("br", 0)
		if err != nil {
			// TODO: print error
			return echo.ErrInternalServerError
		}

		var s sdpS
		err = c.Bind(&s)
		if err != nil {
			return err
		}

		broadcast := &lwebrtc.BroadcasterSession{
			ID: ranId,
		}

		lwebrtc.Broadcasts = append(lwebrtc.Broadcasts, broadcast)

		go lwebrtc.HandleBroadcasterConnection(s.SDP, broadcast, broadcasterLocalSDPChan)

		sdp := <-broadcasterLocalSDPChan

		return c.JSON(http.StatusOK, response{SDP: sdp})
	})

	e.POST("/broadcast/:id", func(c echo.Context) error {
		type broadcastParams struct {
			SDP string `json:"sdp"`
			ID  string `param:"id"`
		}
		type response struct {
			SDP string `json:"sdp"`
		}

		var s broadcastParams
		err := c.Bind(&s)
		if err != nil {
			return err
		}

		fmt.Println("Init endpoint")
		fmt.Printf("Trach Channel has %d data stored", len(localTrackChan))

		var broadcast *lwebrtc.BroadcasterSession

		for _, broadcastSession := range lwebrtc.Broadcasts {
			if broadcastSession.ID == s.ID {
				broadcast = broadcastSession
				break
			}
		}

		if broadcast.Track == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "No broadcast available",
			})
		}

		go lwebrtc.HandleViewer(s.SDP, broadcast.Track, viewerLocalSDPChan)

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
