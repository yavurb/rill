package ui

import (
	"encoding/json"
	"fmt"
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
	routerGroup.GET("/:broadcastID/join", routerCtx.Connect)
}

func (routerCtx *broadcastsRouterCtx) HandleWebsocket(c echo.Context) error {
	ws, err := websocket.Accept(c.Response(), c.Request(), nil)
	if err != nil {
		return HTTPError{Message: "Upgrade to websocket required"}.ErrUpgradeRequired()
	}

	defer ws.Close(websocket.StatusNormalClosure, "goodbye")

	ctx := c.Request().Context()
	for {
		select {
		case <-ctx.Done():
			c.Logger().Info("Request context canceled:", ctx.Err())

			return nil
		default:
			event := new(WsEvent)
			err := wsjson.Read(ctx, ws, event)
			if err != nil {
				if closeStatus := websocket.CloseStatus(err); closeStatus != -1 {
					switch closeStatus {
					case websocket.StatusNormalClosure:
						return nil
					case websocket.StatusGoingAway:
						c.Logger().Info("Client is going away")
						return nil
					case websocket.StatusAbnormalClosure:
						c.Logger().Info("Client is closing abnormally")
						return nil
					}
				}

				c.Logger().Error("Unexpected WebSocket Error:", err)

				return err
			}

			c.Logger().Info("Received: ", event)

			jsonEventData, _ := json.Marshal(event.Data)

			switch event.Event {
			case "subscribe":
				broadcastIn := new(BroadcastIn)

				err = json.Unmarshal(jsonEventData, broadcastIn)
				if err != nil {
					log.Printf("Error: %s", err)
				}

				log.Print(broadcastIn.SDP)
				log.Print(broadcastIn.Title)
			case "unsubscribe":
				fmt.Println("unsubscribe")
			case "broadcast":
				fmt.Println("broadcast")
			default:
				fmt.Println("default")
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

func (routerCtx *broadcastsRouterCtx) CreateBroadcast(c echo.Context) error {
	requestBody := new(BroadcastIn)

	if err := c.Bind(requestBody); err != nil {
		return HTTPError{
			Message: "broadcast sdp and title are required",
		}.ErrUnprocessableEntity()
	}

	if err := c.Validate(requestBody); err != nil {
		return HTTPError{Message: "broadcast sdp and title are required"}.ErrUnprocessableEntity()
	}

	broadcastLocalSDPSession, err := routerCtx.broadcastUsecase.Create(requestBody.SDP, requestBody.Title)
	if err != nil {
		return HTTPError{
			Message: "could no create broadcast",
		}.InternalServerError()
	}

	broadcastOut := &BroadcastCreateOut{
		SDP: broadcastLocalSDPSession,
	}

	return c.JSON(http.StatusCreated, broadcastOut)
}

func (routerCtx *broadcastsRouterCtx) Connect(c echo.Context) error {
	var connectParams BroadcastConnectParams

	if err := c.Bind(&connectParams); err != nil {
		return HTTPError{
			Message: "broadcast sdp is required",
		}.ErrUnprocessableEntity()
	}

	localSDP, err := routerCtx.broadcastUsecase.Connect(connectParams.SDP, connectParams.BroadcastID)
	if err != nil {
		return HTTPError{
			Message: "broadcast not found",
		}.NotFound()
	}

	sdpOut := &BroadcastConnectOut{
		SDP: localSDP,
	}

	return c.JSON(http.StatusOK, sdpOut)
}
