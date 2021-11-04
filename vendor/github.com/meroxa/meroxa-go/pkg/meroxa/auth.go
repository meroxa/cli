package meroxa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cristalhq/jwt/v3"
	"golang.org/x/oauth2"
)

const (
	ClientID = "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb" // TODO this is the CLI ID, create separate client ID for 3rd party apps and provide it as default
)

var (
	OAuth2Endpoint = oauth2.Endpoint{
		AuthURL:  "https://auth.meroxa.io/authorize",
		TokenURL: "https://auth.meroxa.io/oauth/token",
	}
)

// TokenObserver is a function that will be notified when a new token is fetched.
type TokenObserver func(*oauth2.Token)

func DefaultOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID: ClientID,
		Endpoint: OAuth2Endpoint,
	}
}

func newAuthClient(
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
		conf = DefaultOAuth2Config()
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

type tokenSource struct {
	client       *http.Client
	conf         *oauth2.Config
	refreshToken string
	observers    []TokenObserver
}

func (ts *tokenSource) Token() (*oauth2.Token, error) {
	if ts.refreshToken == "" {
		return nil, errors.New("oauth2: token expired and refresh token is not set")
	}

	tk, err := ts.retrieveToken(ts.client, ts.conf, ts.refreshToken)
	if err != nil {
		return nil, err
	}

	if tk.RefreshToken == "" {
		tk.RefreshToken = ts.refreshToken
	}
	if ts.refreshToken != tk.RefreshToken {
		ts.refreshToken = tk.RefreshToken
	}

	// asynchronously notify observers
	go ts.notifyTokenObservers(tk, ts.observers)

	return tk, err
}

func (ts *tokenSource) notifyTokenObservers(token *oauth2.Token, observers []TokenObserver) {
	clone := *token
	for _, o := range observers {
		o(&clone)
	}
}

func (ts *tokenSource) retrieveToken(client *http.Client, conf *oauth2.Config, refreshToken string) (*oauth2.Token, error) {
	tmp := make(map[string]interface{})
	tmp["client_id"] = conf.ClientID
	tmp["grant_type"] = "refresh_token"
	tmp["refresh_token"] = refreshToken
	requestBody, err := json.Marshal(tmp)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(conf.Endpoint.TokenURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}

	if c := resp.StatusCode; c < 200 || c > 299 {
		return nil, &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	// tokenRes is the JSON response body.
	var tokenRes struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"` // relative seconds from now
		// Ignored fields
		// Scope       string `json:"scope"`
		// IDToken     string `json:"id_token"`
	}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	token := &oauth2.Token{
		AccessToken: tokenRes.AccessToken,
		TokenType:   tokenRes.TokenType,
	}
	raw := make(map[string]interface{})
	json.Unmarshal(body, &raw) // no error checks for optional fields
	token = token.WithExtra(raw)

	expiry, err := getTokenExpiry(token.AccessToken)
	if err != nil {
		// fallback to calculate expiry
		expiry = time.Unix(tokenRes.ExpiresIn, 0).UTC()
	}
	token.Expiry = expiry

	// TODO validate JWT token signature
	// keys are available at https://auth.meroxa.io/.well-known/jwks.json

	return token, nil
}

func getTokenExpiry(token string) (time.Time, error) {
	jwtToken, err := jwt.ParseString(token)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse access token as JWT: %w", err)
	}

	var claims jwt.StandardClaims
	err = json.Unmarshal(jwtToken.RawClaims(), &claims)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not parse access token JWT claims: %w", err)
	}

	return claims.ExpiresAt.Time.UTC(), nil
}
