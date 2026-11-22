package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valeryvash/shortener/internal/config"
)

func TestConfigCreation(t *testing.T) {
	testCases := []struct {
		name       string
		args       []string
		wantError  bool
		wantConfig config.Config
	}{
		{
			name: "no commandline args, got default values",
			args: []string{},
			wantConfig: config.Config{
				BindAddress:     "localhost:8080",
				RedirectAddress: "http://localhost:8080",
			},
		},
		{
			name: "provide known args, got not default values",
			args: []string{
				"-a", "localhost:9091",
				"-b", "http://ya.ru:8891",
			},
			wantConfig: config.Config{
				BindAddress:     "localhost:9091",
				RedirectAddress: "http://ya.ru:8891",
			},
		},
		{
			name: "provide unknown args, got error",
			args: []string{
				"-alollo",
			},
			wantError:  true,
			wantConfig: config.Config{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotConfig, err := config.NewConfig(tc.args)

			if tc.wantError {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t,
				tc.wantConfig.BindAddress, gotConfig.BindAddress)
			assert.Equal(t,
				tc.wantConfig.RedirectAddress, gotConfig.RedirectAddress)
		})
	}

}
