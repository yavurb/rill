package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/labstack/echo/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type (
	subscriber struct {
		event chan string
	}
	publisher struct {
		subscribers      map[*subscriber]struct{}
		subscribersMutex sync.Mutex
	}
	broadcastsRouterCtx struct {
		echo             *echo.Echo
		broadcastUsecase domain.BroadcastsUsecases
		publishers       map[string]*publisher
	}
)

func NewBroadcastsRouter(echo *echo.Echo, broadcastUsecase domain.BroadcastsUsecases) {
	routerGroup := echo.Group("broadcasts")
	routerCtx := &broadcastsRouterCtx{
		echo:             echo,
		broadcastUsecase: broadcastUsecase,
		publishers:       make(map[string]*publisher),
	}

	routerGroup.GET("/ws", routerCtx.HandleWebsocket)
	routerGroup.GET("", routerCtx.GetBroadcasts)
	routerGroup.GET("/:id", routerCtx.GetBroadcast)
}

func (routerCtx *broadcastsRouterCtx) HandleWebsocket(c echo.Context) error {
	ws, err := websocket.Accept(
		c.Response(),
		c.Request(),
		&websocket.AcceptOptions{OriginPatterns: []string{
			"localhost:*",
			"rill.one",
			"rill.lat",
		}},
	)
	if err != nil {
		c.Logger().Debug(err)
		return HTTPError{Message: "Upgrade to websocket required"}.ErrUpgradeRequired()
	}

	defer ws.Close(websocket.StatusNormalClosure, "goodbye")

	ctx := c.Request().Context()
	broadcast := new(domain.BroadcastSession)

	for {
		select {
		case <-ctx.Done():
			c.Logger().Info("Request context canceled:", ctx.Err())

			return nil
		case <-broadcast.ContextClose():
			c.Logger().Debug("Broadcast context canceled")
			return nil
		default:
			event := new(WsEvent)
			err := wsjson.Read(ctx, ws, event)
			if err != nil {
				if closeStatus := websocket.CloseStatus(err); closeStatus != -1 {
					switch closeStatus {
					case websocket.StatusNormalClosure:
						broadcast.Close(nil)
						return nil
					case websocket.StatusGoingAway:
						c.Logger().Info("Client is going away")
						broadcast.Close(nil)
						return nil
					case websocket.StatusAbnormalClosure:
						c.Logger().Info("Client is closing abnormally")
						broadcast.Close(err)
						return err
					default:
						c.Logger().Infof("Client closed WebSocket with status: %d", closeStatus)
						broadcast.Close(err)
						return err
					}
				} else if err == io.EOF {
					c.Logger().Info("Client closed the WebSocket connection")
					broadcast.Close(err)
					return err
				}

				c.Logger().Errorf("Unexpected WebSocket Error: %v", err)

				broadcast.Close(err)

				return err
			}

			jsonEventData, _ := json.Marshal(event.Data)

			switch event.Event {
			case "new-broadcast":
				c.Logger().Info("Received new-broadcast event")
				eventData := new(BroadcastIn)

				err = json.Unmarshal(jsonEventData, eventData)
				if err != nil {
					log.Printf("Error: %s", err)
				}

				b, err := routerCtx.broadcastUsecase.Create(eventData.SDP, eventData.Title)
				if err != nil {
					// TODO: Handle the error properly
				}

				broadcastOut := &BroadcastCreateOut{
					SDP: b.LocalSDPSession,
				}
				wsEvent := WsEvent{Event: "new-broadcast", Data: broadcastOut}

				err = wsjson.Write(ctx, ws, wsEvent)
				if err != nil {
					// TODO: Handle the error properly
					return err
				}

				broadcast = b
				defer routerCtx.broadcastUsecase.Delete(broadcast.ID)
			case "new-viewer":
				eventData := new(ViewerIn)

				err = json.Unmarshal(jsonEventData, eventData)
				if err != nil {
					log.Printf("Error: %s", err)
				}

				viewer, err := routerCtx.broadcastUsecase.Connect(eventData.SDP, eventData.BroadcastID)
				if err != nil {
					// TODO: Handle the error properly
				}

				viewerOut := &ViewerOut{
					SDP: viewer.LocalSDPSession,
				}
				wsEvent := WsEvent{Event: "new-viewer", Data: viewerOut}

				err = wsjson.Write(ctx, ws, wsEvent)
				if err != nil {
					// TODO: Handle the error properly
					return err
				}
			case "unsubscribe":
			default:
				// TODO: Handle the case when the event is not recognized. Should we send an error message to the client?
				fmt.Println("No event found")
			}
		}
	}
}

func (routerCtx *broadcastsRouterCtx) GetBroadcasts(c echo.Context) error {
	broadcasts, err := routerCtx.broadcastUsecase.GetBroadcasts()
	if err != nil {
		return HTTPError{
			Message: "unable to get broadcasts",
		}.InternalServerError()
	}

	broadcastsOut := new(BroadcastsOut)

	for _, broascast := range broadcasts {
		broadcastsOut.Broadcasts = append(broadcastsOut.Broadcasts, &BroadcastOut{
			ID:    broascast.ID,
			Title: broascast.Title,
		})
	}

	return c.JSON(http.StatusOK, broadcastsOut)
}

func (routerCtx *broadcastsRouterCtx) GetBroadcast(c echo.Context) error {
	var requestParams GetBroadcastParams

	if err := c.Bind(&requestParams); err != nil {
		return HTTPError{
			Message: "broadcast ID is required",
		}.ErrUnprocessableEntity()
	}

	broadcast, err := routerCtx.broadcastUsecase.Get(requestParams.ID)
	if err != nil {
		return HTTPError{
			Message: "broadcast not found",
		}.NotFound()
	}

	broadcastOut := &BroadcastOut{
		ID:    broadcast.ID,
		Title: broadcast.Title,
	}

	return c.JSON(http.StatusOK, broadcastOut)
}
