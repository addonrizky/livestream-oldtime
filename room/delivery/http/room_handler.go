package http

import (
	"context"
	"net/http"

	"github.com/asumsi/livestream/middleware"
	"github.com/asumsi/livestream/modelgeneral"
	room "github.com/asumsi/livestream/room"
	"github.com/asumsi/livestream/room/models"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

var validate *validator.Validate

type HttpRoomHandler struct {
	AUsecase room.RoomUsecase
}

func NewRoomHttpHandler(e *echo.Echo, as room.RoomUsecase) {
	handler := &HttpRoomHandler{
		AUsecase: as,
	}

	routeGroup := e.Group("/api/v1")
	// Middleware Auth
	routeGroup.Use(echo.WrapMiddleware(middleware.RequestLogMiddleware))
	routeGroup.Use(echo.WrapMiddleware(middleware.JWTAuthorizationMiddleware))
	routeGroup.Use(middleware.NewAppContextMiddleware())

	routeGroup.POST("/rooms", handler.GetRooms)

}

func (a *HttpRoomHandler) GetRooms(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req modelgeneral.Pagination
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, modelgeneral.LoginResponse{Code: http.StatusUnprocessableEntity, Message: err.Error()})
	}

	data, err := a.AUsecase.GetRooms(ctx, req.DataPerPage, req.Page)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, models.RoomResponse{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Result: data})
}
