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
			ID: broascast.ID,
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
		ID: broadcast.ID,
	}

	return c.JSON(http.StatusOK, broadcastOut)
}

func (routerCtx *broadcastsRouterCtx) CreateBroadcast(c echo.Context) error {
	var requestBody BroadcastIn

	if err := c.Bind(&requestBody); err != nil {
		return HTTPError{
			Message: "broadcast sdp is required",
		}.ErrUnprocessableEntity()
	}

	broadcastLocalSDPSession, err := routerCtx.broadcastUsecase.Create(requestBody.SDP)
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
