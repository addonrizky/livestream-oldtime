package http

import (
	"net/http"

	"github.com/asumsi/livestream/modelgeneral"
	"github.com/asumsi/livestream/stream/models"
	"github.com/asumsi/livestream/stream/usecase"
	"github.com/labstack/echo/v4"
)

type httpStreamHandler struct {
	streamUsecase usecase.StreamUsecase
}

func NewStreamHttpHandler(e *echo.Echo, suc usecase.StreamUsecase) {
	handler := &httpStreamHandler{
		streamUsecase: suc,
	}

	e.POST("api/v1/stream/init_session", handler.initSession)
	e.POST("api/v1/stream/start_streaming", handler.startStreaming)
	e.POST("api/v1/stream/end_streaming", handler.endStreaming)
}

// InitSession godoc
// @Summary      Init session stream
// @Description  Init session send offer session to WebRTC
// @Tags         stream
// @Accept       json
// @Produce      json
// @Param 		 body body models.InitSessionRequest false "Body Request"
// @Success      200  {object}  models.ResponseSwag
// @Router       /stream/init_session [post]
func (a *httpStreamHandler) initSession(c echo.Context) error {
	var err error
	var errUsecase *models.ErrorObject
	var resultUsecase map[string]interface{}
	u := new(models.InitSessionRequest)
	response := modelgeneral.NewResponse("oihfgh02hg0932h9")

	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(u); err != nil {
		return c.JSON(http.StatusOK, "validation error")
	}

	resultUsecase, errUsecase = a.streamUsecase.InitSession(u.OfferSession, u.RoomId, u.UserId)

	if errUsecase != nil {
		response.Code = errUsecase.ErrorCode
		response.Desc = errUsecase.ErrorMessage
		return c.JSON(http.StatusOK, response)
	}

	response.Code = "00"
	response.Desc = "Transaksi Sukses"
	response.Data = resultUsecase
	return c.JSON(http.StatusOK, response)
}

func (a *httpStreamHandler) startStreaming(c echo.Context) error {
	var err error

	u := new(models.StartSessionRequest)
	response := modelgeneral.NewResponse("oihfgh02hg0932h9")

	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(u); err != nil {
		return c.JSON(http.StatusOK, "validation error")
	}

	resultUsecase, _ := a.streamUsecase.StartStreaming(u.RoomId)

	response.Code = "00"
	response.Desc = "Streaming Run"
	response.Data = resultUsecase

	return c.JSON(http.StatusOK, response)
}

func (a *httpStreamHandler) endStreaming(c echo.Context) error {
	var err error

	u := new(models.EndSessionRequest)
	response := modelgeneral.NewResponse("oihfgh02hg0932h9")

	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(u); err != nil {
		return c.JSON(http.StatusOK, "validation error")
	}

	resultUsecase, errUsecase := a.streamUsecase.EndStreaming(u.PID)

	if errUsecase != nil {
		response.Code = errUsecase.ErrorCode
		response.Desc = errUsecase.ErrorMessage
		return c.JSON(http.StatusOK, response)
	}

	response.Code = "00"
	response.Desc = "Streaming Stop"
	response.Data = resultUsecase

	return c.JSON(http.StatusOK, response)
}
