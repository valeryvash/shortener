package handler_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/valeryvash/shortener/internal/handler"
	hMocks "github.com/valeryvash/shortener/internal/handler/mocks"
	"github.com/valeryvash/shortener/internal/service"
)

var (
	_               handler.Service = (*service.Service)(nil)
	redirectAddress                 = "http://localhost:8080"
)

func TestCreateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type want struct {
		code        int
		response    string
		contentType string
	}

	tcs := []struct {
		name        string
		method      string
		requestBody string
		service     func(m *hMocks.Service)
		want        want
	}{
		{
			name:        "succesfully created",
			method:      http.MethodPost,
			requestBody: "http://redirect.me.here",
			service: func(m *hMocks.Service) {
				m.EXPECT().CreateShortAddr("http://redirect.me.here").Times(1).Return("AbCd1234", nil)
			},
			want: want{
				code:        201,
				response:    fmt.Sprintf("%s/%s", redirectAddress, "AbCd1234"),
				contentType: "text/plain",
			},
		},
		{
			name:        "error from service",
			method:      http.MethodPost,
			requestBody: "http://provide.an.error",
			service: func(m *hMocks.Service) {
				m.EXPECT().CreateShortAddr("http://provide.an.error").Times(1).Return("", errors.New("service error"))
			},
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.requestBody))
			w := httptest.NewRecorder()

			mockService := hMocks.NewService(t)
			tc.service(mockService)

			h := handler.NewHandler(mockService, redirectAddress)

			e := gin.New()
			h.RegisterRoutesGIN(e)

			e.ServeHTTP(w, r)

			body, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			require.Equal(t, tc.want.code, w.Result().StatusCode)

			require.True(t,
				strings.Contains(w.Result().Header.Get("Content-type"), tc.want.contentType),
			)

			require.Equal(t, tc.want.response, string(body))
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type want struct {
		code     int
		response string
		location string
	}

	tcs := []struct {
		name       string
		method     string
		requestURL string
		service    func(m *hMocks.Service)
		want       want
	}{
		{
			name:       "succesfully redirected",
			method:     http.MethodGet,
			requestURL: "/12345678",
			service: func(m *hMocks.Service) {
				m.EXPECT().GetFullAddr("12345678").Times(1).Return("http://wanted.redirect.ad", nil)
			},
			want: want{
				code:     307,
				location: "http://wanted.redirect.ad",
			},
		},
		{
			name:       "error: short url",
			method:     http.MethodGet,
			requestURL: "/1234567",
			service:    func(m *hMocks.Service) {},
			want: want{
				code:     400,
				location: "",
			},
		},
		{
			name:       "error from service",
			method:     http.MethodGet,
			requestURL: "/12345678",
			service: func(m *hMocks.Service) {
				m.EXPECT().GetFullAddr("12345678").Times(1).Return("", errors.New("error"))
			},
			want: want{
				code:     400,
				location: "",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.requestURL, nil)
			w := httptest.NewRecorder()

			mockService := hMocks.NewService(t)
			tc.service(mockService)

			h := handler.NewHandler(mockService, redirectAddress)

			e := gin.New()
			h.RegisterRoutesGIN(e)

			e.ServeHTTP(w, r)

			require.Equal(t, tc.want.code, w.Result().StatusCode)

			require.Equal(t, tc.want.location, w.Result().Header.Get("Location"))
		})
	}
}
