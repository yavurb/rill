// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package signal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTTPSDPServer starts a HTTP Server that consumes SDPs
func HTTPSDPServer(port int, localSDPChan chan string) chan string {
	sdpChan := make(chan string)
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4321"},
	}))

	e.POST("/", func(c echo.Context) error {
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

		sdpChan <- string(s.SDP)

		sdp := <-localSDPChan

		fmt.Printf("Got SDP, %s \n", sdp)

		return c.JSON(http.StatusOK, response{SDP: sdp})
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	body, _ := io.ReadAll(r.Body)
	// 	sdpChan <- string(body)

	// 	sdp := <-localSDPChan

	// 	fmt.Printf("Got SDP, %s \n", sdp)
	// 	fmt.Fprintf(w, sdp)
	// })

	go func() {
		// nolint: gosec
		e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
		// err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		// if err != nil {
		// 	panic(err)
		// }
	}()

	return sdpChan
}
