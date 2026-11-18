package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/valeryvash/shortener/internal/handler"
	hMocks "github.com/valeryvash/shortener/internal/handler/mocks"
	"github.com/valeryvash/shortener/internal/service"
)

var _ handler.Service = (*service.Service)(nil)

func TestCreateHandler(t *testing.T) {

	type want struct {
		code        int
		response    string
		contentType string
	}

	tcs := []struct {
		name        string
		method      string
		requestBody string
		service     func() *hMocks.Service
		want        want
	}{
		{
			name:        "succesfully created",
			method:      http.MethodPost,
			requestBody: "text",
			service: func() *hMocks.Service {
				m := hMocks.Service{}
				m.EXPECT().CreateShortAddr("text").Times(1).Return("abc123", nil)
				m.EXPECT().GetFullAddr(mock.Anything).Times(0).Return("", nil)
				return &m
			},
			want: want{
				code:        201,
				response:    "http://localhost:8080/abc123",
				contentType: "text/plain",
			},
		},
		{
			name:        "error from service",
			method:      http.MethodPost,
			requestBody: "text",
			service: func() *hMocks.Service {
				m := hMocks.Service{}
				m.EXPECT().CreateShortAddr(mock.Anything).Times(1).Return("", errors.New("service error"))
				m.EXPECT().GetFullAddr(mock.Anything).Times(0).Return("", nil)
				return &m
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

			h := handler.NewHandler(tc.service())

			h.CreateHandler(w, r)

			body, err := io.ReadAll(w.Result().Body)
			defer w.Result().Body.Close()

			require.NoError(t, err)

			require.Equal(t, tc.want.code, w.Result().StatusCode)

			require.Equal(t, tc.want.contentType, w.Result().Header.Get("Content-type"))

			require.Equal(t, tc.want.response, string(body))
		})
	}
}

func TestRedirectHandler(t *testing.T) {

	type want struct {
		code     int
		response string
		location string
	}

	tcs := []struct {
		name       string
		method     string
		requestURL string
		service    func() *hMocks.Service
		want       want
	}{
		{
			name:       "succesfully redirected",
			method:     http.MethodGet,
			requestURL: "/12345678",
			service: func() *hMocks.Service {
				m := hMocks.Service{}
				m.EXPECT().CreateShortAddr(mock.Anything).Times(0).Return("", nil)
				m.EXPECT().GetFullAddr("12345678").Times(1).Return("http://ya.ru", nil)
				return &m
			},
			want: want{
				code:     307,
				location: "http://ya.ru",
			},
		},
		{
			name:       "error from service",
			method:     http.MethodGet,
			requestURL: "/12345678",
			service: func() *hMocks.Service {
				m := hMocks.Service{}
				m.EXPECT().CreateShortAddr(mock.Anything).Times(0).Return("", nil)
				m.EXPECT().GetFullAddr(mock.Anything).Times(1).Return("", errors.New("error"))
				return &m
			},
			want: want{
				code:     400,
				location: "",
			},
		},
		{
			name:       "error: short url",
			method:     http.MethodGet,
			requestURL: "/1234567",
			service: func() *hMocks.Service {
				m := hMocks.Service{}
				m.EXPECT().CreateShortAddr(mock.Anything).Times(0).Return("", nil)
				m.EXPECT().GetFullAddr(mock.Anything).Times(0).Return("", nil)
				return &m
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

			h := handler.NewHandler(tc.service())

			h.RedirectHandler(w, r)

			require.Equal(t, tc.want.code, w.Result().StatusCode)

			require.Equal(t, tc.want.location, w.Result().Header.Get("Location"))
		})
	}
}
