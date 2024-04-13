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
}

func (routerCtx *broadcastsRouterCtx) GetBroadcast(c echo.Context) error {
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

func (routerCtx *broadcastsRouterCtx) GetBroadcasts(c echo.Context) error {
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
