package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type broadcastsRouterCtx struct {
	echo             *echo.Echo
	broadcastUsecase domain.BroadcastsUsecases
}

func NewBroadcastsRouter(echo *echo.Echo, broadcastUsecase domain.BroadcastsUsecases) {
	routerGroup := echo.Group("broadcasts")
	routerCtx := &broadcastsRouterCtx{
		echo,
		broadcastUsecase,
	}

	routerGroup.GET("", routerCtx.GetBroadcasts)
	routerGroup.GET("/:id", routerCtx.GetBroadcast)
	routerGroup.POST("", routerCtx.CreateBroadcast)
	routerGroup.POST("/:broadcastID/join", routerCtx.Connect)
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
