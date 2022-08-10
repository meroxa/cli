package platform

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"golang.org/x/oauth2"
)

type ClientConfig struct {
	AccessToken  string `env:"ACCESS_TOKEN"`
	RefreshToken string `env:"REFRESH_TOKEN"`
	AuthAudience string `env:"MEROXA_AUTH_AUDIENCE" envDefault:"https://api.meroxa.io/v1"`
	AuthDomain   string `env:"MEROXA_AUTH_DOMAIN" envDefault:"auth.meroxa.io"`
	AuthClientID string `env:"MEROXA_AUTH_CLIENT_ID" envDefault:"2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb"`
}

const Version = "0.1.0"

var cfg ClientConfig

type Client struct {
	meroxa.Client
}

func newClient() (*Client, error) {
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	options := []meroxa.Option{
		meroxa.WithUserAgent(fmt.Sprintf("Meroxa CLI %s", Version)),
	}

	if overrideAPIURL := os.Getenv("API_URL"); overrideAPIURL != "" {
		options = append(options, meroxa.WithBaseURL(overrideAPIURL))
	} else if overrideAPIURL := os.Getenv("MEROXA_API_URL"); overrideAPIURL != "" {
		options = append(options, meroxa.WithBaseURL(overrideAPIURL))
	}

	if os.Getenv("MEROXA_DEBUG") != "" {
		options = append(options, meroxa.WithDumpTransport(log.Writer()))
	}

	options = append(options, meroxa.WithAuthentication(
		&oauth2.Config{
			ClientID: cfg.AuthClientID,
			Endpoint: oauthEndpoint(cfg.AuthDomain, cfg.AuthAudience),
		},
		cfg.AccessToken,
		cfg.RefreshToken,
		onTokenRefreshed,
	))

	mc, err := meroxa.New(options...)
	return &Client{mc}, err
}

func oauthEndpoint(domain, audience string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("https://%s/authorize?audience=%s", domain, audience),
		TokenURL: fmt.Sprintf("https://%s/oauth/token", domain),
	}
}

func onTokenRefreshed(token *oauth2.Token) {
	cfg.AccessToken = token.AccessToken
	cfg.RefreshToken = token.RefreshToken
}
