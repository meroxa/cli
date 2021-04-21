package auth

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	ClientID = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb" // TODO this is the CLI ID, create separate client ID for 3rd party apps and provide it as default
)

var (
	Endpoint = oauth2.Endpoint{
		AuthURL:  "https://auth.meroxa.io/authorize",
		TokenURL: "https://auth.meroxa.io/oauth/token",
	}
)

// TokenObserver is a function that will be notified when a new token is fetched.
type TokenObserver func(*oauth2.Token)

func DefaultConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID: ClientID,
		Endpoint: Endpoint,
	}
}

func NewClient(
	client *http.Client,
	conf *oauth2.Config,
	accessToken,
	refreshToken string,
	tokenObservers ...TokenObserver,
) (*http.Client, error) {
	if client == nil {
		client = http.DefaultClient
	}
	if conf == nil {
		conf = DefaultConfig()
	}

	var expiry time.Time
	if accessToken != "" {
		var err error
		expiry, err = getTokenExpiry(accessToken)
		if err != nil {
			return nil, err
		}
	}

	ts := oauth2.ReuseTokenSource(
		&oauth2.Token{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			RefreshToken: refreshToken,
			Expiry:       expiry,
		},
		&tokenSource{
			conf:         conf,
			client:       client,
			refreshToken: refreshToken,
			observers:    tokenObservers,
		},
	)

	// make sure the oauth2 client is using the supplied client for outgoing requests
	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, client)

	return oauth2.NewClient(ctx, ts), nil
}
