package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valeryvash/shortener/internal/handler"
)

type serviceMock struct {
	createShortAddr func(fullAddr string) (shortAddr string, err error)
	getFullAddr     func(shortAddr string) (fullAddr string, err error)

	CreateShortAddrArgs []string
	GetFullAddrArgs     []string
}

func (s *serviceMock) CreateShortAddr(fullAddr string) (shortAddr string, err error) {
	s.CreateShortAddrArgs = append(s.CreateShortAddrArgs, fullAddr)
	return s.createShortAddr(fullAddr)
}

func (s *serviceMock) GetFullAddr(shortAddr string) (fullAddr string, err error) {
	s.GetFullAddrArgs = append(s.GetFullAddrArgs, shortAddr)
	return s.getFullAddr(fullAddr)
}

var _ handler.Service = (*serviceMock)(nil)

func TestCreateHandler(t *testing.T) {

	type want struct {
		code            int
		response        string
		contentType     string
		createShortArgs []string
		getFullAddrArgs []string
	}

	tcs := []struct {
		name        string
		method      string
		requestBody string
		service     *serviceMock
		want        want
	}{
		{
			name:        "succesfully created",
			method:      http.MethodPost,
			requestBody: "text",
			service: &serviceMock{
				createShortAddr: func(fullAddr string) (shortAddr string, err error) { return "abc123", nil },
				getFullAddr:     func(shortAddr string) (fullAddr string, err error) { return "", nil },

				CreateShortAddrArgs: []string{},
				GetFullAddrArgs:     []string{},
			},
			want: want{
				code:            201,
				response:        "http://localhost:8080/abc123",
				contentType:     "text/plain",
				createShortArgs: []string{"text"},
				getFullAddrArgs: []string{},
			},
		},
		{
			name:        "error from service",
			method:      http.MethodPost,
			requestBody: "text",
			service: &serviceMock{
				createShortAddr: func(fullAddr string) (shortAddr string, err error) { return "", errors.New("error") },
				getFullAddr:     func(shortAddr string) (fullAddr string, err error) { return "", nil },

				CreateShortAddrArgs: []string{},
				GetFullAddrArgs:     []string{},
			},
			want: want{
				code:            400,
				response:        "",
				contentType:     "",
				createShortArgs: []string{"text"},
				getFullAddrArgs: []string{},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.requestBody))
			w := httptest.NewRecorder()

			h := handler.NewHandler(tc.service)

			h.CreateHandler(w, r)

			body, err := io.ReadAll(w.Result().Body)
			defer w.Result().Body.Close()

			require.NoError(t, err)

			require.Equal(t, tc.want.code, w.Result().StatusCode)

			require.Equal(t, tc.want.contentType, w.Result().Header.Get("Content-type"))

			require.Equal(t, tc.want.response, string(body))

			assert.Equal(t, tc.want.createShortArgs, tc.service.CreateShortAddrArgs)
			assert.Equal(t, tc.want.getFullAddrArgs, tc.service.GetFullAddrArgs)
		})
	}
}

func TestRedirectHandler(t *testing.T) {

	type want struct {
		code            int
		response        string
		location        string
		createShortArgs []string
		getFullAddrArgs []string
	}

	tcs := []struct {
		name       string
		method     string
		requestURL string
		service    *serviceMock
		want       want
	}{
		{
			name:       "succesfully redirected",
			method:     http.MethodGet,
			requestURL: "/12345678",
			service: &serviceMock{
				createShortAddr: func(fullAddr string) (shortAddr string, err error) { return "", nil },
				getFullAddr:     func(shortAddr string) (fullAddr string, err error) { return "http://ya.ru", nil },

				CreateShortAddrArgs: []string{},
				GetFullAddrArgs:     []string{},
			},
			want: want{
				code:            307,
				location:        "http://ya.ru",
				createShortArgs: []string{},
				getFullAddrArgs: []string{"12345678"},
			},
		},
		{
			name:       "error from service",
			method:     http.MethodGet,
			requestURL: "/12345678",
			service: &serviceMock{
				createShortAddr: func(fullAddr string) (shortAddr string, err error) { return "", nil },
				getFullAddr:     func(shortAddr string) (fullAddr string, err error) { return "", errors.New("error") },

				CreateShortAddrArgs: []string{},
				GetFullAddrArgs:     []string{},
			},
			want: want{
				code:            400,
				location:        "",
				createShortArgs: []string{},
				getFullAddrArgs: []string{"12345678"},
			},
		},
		{
			name:       "error: short url",
			method:     http.MethodGet,
			requestURL: "/1234567",
			service: &serviceMock{
				createShortAddr: func(fullAddr string) (shortAddr string, err error) { return "", nil },
				getFullAddr:     func(shortAddr string) (fullAddr string, err error) { return "", nil },

				CreateShortAddrArgs: []string{},
				GetFullAddrArgs:     []string{},
			},
			want: want{
				code:            400,
				location:        "",
				createShortArgs: []string{},
				getFullAddrArgs: []string{},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.requestURL, nil)
			w := httptest.NewRecorder()

			h := handler.NewHandler(tc.service)

			h.RedirectHandler(w, r)

			require.Equal(t, tc.want.code, w.Result().StatusCode)

			require.Equal(t, tc.want.location, w.Result().Header.Get("Location"))

			assert.Equal(t, tc.want.createShortArgs, tc.service.CreateShortAddrArgs)
			assert.Equal(t, tc.want.getFullAddrArgs, tc.service.GetFullAddrArgs)
		})
	}
}
