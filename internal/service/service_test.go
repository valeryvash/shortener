package service_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/valeryvash/shortener/internal/handler"
	"github.com/valeryvash/shortener/internal/service"
	sMock "github.com/valeryvash/shortener/internal/service/mocks"
)

var _ handler.Service = (*service.Service)(nil)

func TestCreateShortAddr(t *testing.T) {

	testCases := []struct {
		name      string
		repo      func(m *sMock.Repository)
		fullAddr  string
		wantError bool
	}{
		{
			name:     "succesfully created short address",
			fullAddr: "http://give.me.shorten/link",
			repo: func(m *sMock.Repository) {
				m.EXPECT().Save("http://give.me.shorten/link", mock.Anything).Times(1).Return(nil)
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := sMock.NewRepository(t)
			tc.repo(mockRepo)
			service := service.NewService(mockRepo)

			shortValue, err := service.CreateShortAddr(tc.fullAddr)

			if !tc.wantError {
				require.NoError(t, err)
				require.NotEmpty(t, shortValue)
			}
		})
	}
}
