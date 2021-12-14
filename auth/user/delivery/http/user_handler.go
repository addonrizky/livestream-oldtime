package http

import (
	"context"
	"fmt"
	"net/http"

	user "github.com/asumsi/livestream/auth/user"
	"github.com/asumsi/livestream/auth/user/models"
	"github.com/asumsi/livestream/modelgeneral"
	"github.com/asumsi/livestream/utility"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type HttpUserHandler struct {
	userUsecase user.UserUsecase
}

func NewUserHttpHandler(e *echo.Echo, userUsecase user.UserUsecase) {
	handler := &HttpUserHandler{
		userUsecase: userUsecase,
	}

	e.POST("/api/v1/register", handler.Register)
	e.POST("/api/v1/login", handler.Login)
	e.POST("/api/v1/send-reset-password", handler.SendEmailResetPassword)
	e.POST("/api/v1/reset-password", handler.ResetPassword)
}

func (a *HttpUserHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var userData models.UserReq
	err := c.Bind(&userData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{Code: http.StatusUnprocessableEntity, Message: err.Error(), Result: false})
	}

	userData.IsActive = true

	if ok, err := isRequestValid(&userData); !ok {
		return c.JSON(http.StatusBadRequest, models.Response{Code: http.StatusBadRequest, Message: err.Error(), Result: false})
	}

	err = a.userUsecase.CreateUser(ctx, &userData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Code: http.StatusBadRequest, Message: err.Error(), Result: false})
	}
	return c.JSON(http.StatusCreated, models.Response{Code: http.StatusCreated, Message: http.StatusText(http.StatusCreated), Result: true})
}

func (a *HttpUserHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req modelgeneral.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, modelgeneral.LoginResponse{Code: http.StatusUnprocessableEntity, Message: err.Error()})
	}

	data, err := a.userUsecase.Authenticate(ctx, req, "APP")
	if err != nil {
		return c.JSON(http.StatusBadRequest, modelgeneral.LoginResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}

	result := modelgeneral.LoginAttribute{
		Token:        data.Token,
		RefreshToken: data.RefreshToken,
	}

	return c.JSON(http.StatusCreated, modelgeneral.LoginResponse{Code: http.StatusCreated, Message: http.StatusText(http.StatusCreated), Result: result})
}

func isRequestValid(m *models.UserReq) (bool, error) {
	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *HttpUserHandler) SendEmailResetPassword(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req models.ResetPasswordReq
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.EmailResetPasswordResponse{Code: http.StatusUnprocessableEntity, Message: err.Error(), ResetPasswordLink: ""})
	}

	encryptEmail, err := utility.Encrypter(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.EmailResetPasswordResponse{Code: http.StatusBadRequest, Message: err.Error(), ResetPasswordLink: ""})
	}

	// send email via SMTP

	prefixLink := utility.GetConfigString("url")
	resetPasswordLink := prefixLink + "/api/v1/reset-password?key=" + encryptEmail

	return c.JSON(http.StatusOK, models.EmailResetPasswordResponse{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), ResetPasswordLink: resetPasswordLink})
}

func (a *HttpUserHandler) ResetPassword(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var userData models.ResetPasswordReq
	err := c.Bind(&userData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{Code: http.StatusUnprocessableEntity, Message: err.Error(), Result: false})
	}

	key := c.QueryParam("key")
	emailUser, err := utility.Decrypter(key)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Code: http.StatusBadRequest, Message: err.Error(), Result: false})
	}

	userData.Email = emailUser

	err = a.userUsecase.UpdatePassword(ctx, &userData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Code: http.StatusBadRequest, Message: err.Error(), Result: false})
	}

	return c.JSON(http.StatusOK, models.Response{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Result: true})

}

func (a *HttpUserHandler) EncryptTest(c echo.Context) error {

	test, err := utility.Encrypter("hendrik@asumsi.co")
	fmt.Println("test: ", test)
	fmt.Println("err: ", err)

	dec, err := utility.Decrypter("45754beef5347334ee63d6d8.41b181db9071f5651a4ec669d75d62d3ac009d580f677db0d3323946560e5542f1")
	fmt.Println("dec: ", dec)
	fmt.Println("err: ", err)
	return c.JSON(http.StatusOK, models.Response{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Result: true})
}
