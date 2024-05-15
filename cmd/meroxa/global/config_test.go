package global

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetMeroxaMeroxaAuthCallbackURL(t *testing.T) {
	oldConfig := Config

	testCases := []struct {
		name   string
		config func() *viper.Viper
		want   string
	}{
		{
			name: "when MEROXA_AUTH_CALLBACK_URL is not set",
			config: func() *viper.Viper {
				cfg := viper.New()
				return cfg
			},
			want: "http://localhost:21900/oauth/callback",
		},
		{
			name: "when MEROXA_AUTH_CALLBACK_URL is set",
			config: func() *viper.Viper {
				cfg := viper.New()
				cfg.Set(MeroxaAuthCallbackURL, "https://nimbus.meroxa.io:3000/oauth/callback")
				return cfg
			},
			want: "https://nimbus.meroxa.io:3000/oauth/callback",
		},
		{
			name: "when MEROXA_AUTH_CALLBACK_PROTOCOL is set",
			config: func() *viper.Viper {
				cfg := viper.New()
				cfg.Set(MeroxaAuthCallbackProtocol, "https")
				return cfg
			},
			want: "https://localhost:21900/oauth/callback",
		},
		{
			name: "when MEROXA_AUTH_CALLBACK_HOST is set",
			config: func() *viper.Viper {
				cfg := viper.New()
				cfg.Set(MeroxaAuthCallbackHost, "nimbus.meroxa.io:3000")
				return cfg
			},
			want: "http://nimbus.meroxa.io:3000/oauth/callback",
		},
		{
			name: "when MEROXA_AUTH_CALLBACK_PORT is set",
			config: func() *viper.Viper {
				cfg := viper.New()
				cfg.Set(MeroxaAuthCallbackPort, "3000")
				return cfg
			},
			want: "http://localhost:3000/oauth/callback",
		},
		{
			name: "when MEROXA_AUTH_CALLBACK_URL, MEROXA_AUTH_CALLBACK_PROTOCOL, and MEROXA_AUTH_CALLBACK_HOST are set",
			config: func() *viper.Viper {
				cfg := viper.New()
				cfg.Set(MeroxaAuthCallbackURL, "https://nimbus.meroxa.io:3000/oauth/callback")
				cfg.Set(MeroxaAuthCallbackPort, "https")
				cfg.Set(MeroxaAuthCallbackHost, "nimbus.meroxa.io")

				return cfg
			},
			want: "https://nimbus.meroxa.io:3000/oauth/callback",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Config = tc.config()
			got := GetMeroxaAuthCallbackURL()

			if got != tc.want {
				t.Errorf("expected MEROXA_AUTH_CALLBACK_URL to be %q, got %q", tc.want, got)
			}
		})
	}

	Config = oldConfig
}
