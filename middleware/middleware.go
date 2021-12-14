package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/asumsi/livestream/utility"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type (
	AppContext struct {
		echo.Context
	}
	ErrorModel struct {
		Error string
	}
)

func NewAppContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				appContext AppContext
				ok         bool
			)

			if appContext, ok = c.(AppContext); !ok {
				appContext = AppContext{
					Context: c,
				}
			}

			return next(appContext)
		}
	}
}

func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		reqID := fmt.Sprintf("%v", uuid.New())
		r.Header.Set("request_id", reqID)

		b := bytes.NewBuffer(make([]byte, 0))

		reader := io.TeeReader(r.Body, b)
		headerByt, _ := json.Marshal(r.Header)
		bodyBytStr := getBodyByteStr(reader)

		r.Body = ioutil.NopCloser(b)

		fmt.Println(nil, "[IN_REQUEST: ", r.URL, "] HEADER:", string(headerByt), " BODY:", bodyBytStr)
		next.ServeHTTP(w, r)
	})
}

// skip log for base64
func getBodyByteStr(reader io.Reader) string {

	bodyByt := []byte{}
	var req interface{}
	json.NewDecoder(reader).Decode(&req)
	bodyByt, _ = json.Marshal(req)

	return string(bodyByt)
}

func JWTAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			jsonError(w, ErrorModel{Error: "Invalid Token"}, http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		method := utility.GetConfigString(`security.jwt.method`)
		secret := utility.GetConfigString(`security.jwt.secret`)
		token, err := utility.ValidateJwtToken(tokenString, method, secret)
		if token != nil && err == nil {
			next.ServeHTTP(w, r)
		} else {
			jsonError(w, ErrorModel{Error: "Invalid Token"}, http.StatusUnauthorized)
		}
	})
}

func jsonError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
