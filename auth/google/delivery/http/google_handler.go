package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	googleAuth "github.com/asumsi/livestream/auth/google"
	"github.com/asumsi/livestream/auth/google/models"
	"github.com/asumsi/livestream/auth/user"
	modelUser "github.com/asumsi/livestream/auth/user/models"
	modelgeneral "github.com/asumsi/livestream/modelgeneral"
	"github.com/asumsi/livestream/utility"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type HttpGoogleHandler struct {
	googleUsecase googleAuth.GoogleUsecase
	userUsecase   user.UserUsecase
}

func NewGoogleHttpHandler(e *echo.Echo, googleUsecase googleAuth.GoogleUsecase, userUsecase user.UserUsecase) {
	handler := &HttpGoogleHandler{
		googleUsecase: googleUsecase,
		userUsecase:   userUsecase,
	}

	e.POST("/api/v1/google-login", handler.GoogleLogin)
	e.GET("/api/v1/google-callback", handler.GoogleCallback)
}

func (a *HttpGoogleHandler) GoogleLogin(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		oauthConf = &oauth2.Config{
			ClientID:     utility.GetConfigString(`google_clientId`),
			ClientSecret: utility.GetConfigString(`google_clientSecret`),
			RedirectURL:  utility.GetConfigString(`google_redirectURL`),
			Scopes:       []string{utility.GetConfigString(`google_scopeURL`)},
			Endpoint:     google.Endpoint,
		}
		oauthStateString = utility.GetConfigString("google_oauth_StateString")
	)

	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: err.Error(), Result: ""})
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()

	return c.JSON(http.StatusOK, models.GoogleResponse{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Result: url})
}

func (a *HttpGoogleHandler) GoogleCallback(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		oauthConf = &oauth2.Config{
			ClientID:     utility.GetConfigString(`google_clientId`),
			ClientSecret: utility.GetConfigString(`google_clientSecret`),
			RedirectURL:  utility.GetConfigString(`google_redirectURL`),
			Scopes:       []string{utility.GetConfigString(`google_scopeURL`)},
			Endpoint:     google.Endpoint,
		}
		oauthStateString = utility.GetConfigString("google_oauth_StateString")
	)

	state := c.QueryParam("state")
	if state != oauthStateString {
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: "unauthorized", Result: ""})
	}

	code := c.QueryParam("code")
	if code == "" {
		reason := c.QueryParam("error_reason")
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: reason, Result: ""})
	}

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: err.Error(), Result: ""})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: err.Error(), Result: ""})
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: err.Error(), Result: ""})
	}

	res := models.ResponseFromGoogle{}
	json.Unmarshal([]byte(response), &res)

	check, err := a.googleUsecase.CheckIsUserAlreadyRegistered(ctx, res.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: err.Error(), Result: ""})
	}

	if !check {
		var userData modelUser.UserReq
		userData.Email = res.Email
		userData.Name = res.Email
		userData.LoginType = "GOOGLE"
		userData.RoleID = int64(1)
		userData.IsActive = true

		err = a.userUsecase.CreateUser(ctx, &userData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: err.Error(), Result: ""})
		}
	}

	var req modelgeneral.LoginRequest
	req.Email = res.Email
	req.Username = res.Email
	data, err := a.userUsecase.Authenticate(ctx, req, "GOOGLE")
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.GoogleResponse{Code: http.StatusBadRequest, Message: err.Error(), Result: ""})
	}

	accessToken := data.Token
	refreshToken := data.RefreshToken

	result := modelgeneral.LoginAttribute{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, modelgeneral.LoginResponse{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Result: result})
}
