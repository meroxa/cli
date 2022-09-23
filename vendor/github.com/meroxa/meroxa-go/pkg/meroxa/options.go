package meroxa

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

type Option func(*Requester) error

// WithBaseURL sets the base url in the client.
// The default is "https://api.meroxa.io".
func WithBaseURL(rawurl string) Option {
	return func(r *Requester) error {
		u, err := url.Parse(rawurl)
		if err != nil {
			return err
		}
		r.baseURL = u
		return nil
	}
}

// WithClientTimeout sets the http client timeout.
// The default is 5 seconds.
func WithClientTimeout(timeout time.Duration) Option {
	return func(r *Requester) error {
		r.httpClient.Timeout = timeout
		return nil
	}
}

// WithUserAgent sets the User-Agent header.
// The default is "meroxa-go".
func WithUserAgent(ua string) Option {
	return func(r *Requester) error {
		r.userAgent = ua
		return nil
	}
}

// WithDumpTransport will dump the outgoing requests and incoming responses and
// write them to writer.
func WithDumpTransport(writer io.Writer) Option {
	return func(r *Requester) error {
		r.httpClient.Transport = &dumpTransport{
			out:                    writer,
			transport:              r.httpClient.Transport,
			obfuscateAuthorization: true,
		}
		return nil
	}
}

// WithClient sets the http client to use for requests.
func WithClient(httpClient *http.Client) Option {
	return func(r *Requester) error {
		r.httpClient = httpClient
		return nil
	}
}

// WithAuthentication sets an authenticated http client that takes care of
// adding the access token to requests as well as refreshing it with the
// refresh token when it expires. Observers will be called each time the token
// is refreshed.
// Note: provide WithClientTimeout option before WithAuthentication to set the
// timeout of the client used for fetching access tokens.
func WithAuthentication(conf *oauth2.Config, accessToken, refreshToken string, observers ...TokenObserver) Option {
	return func(r *Requester) error {
		httpClient, err := newAuthClient(r.httpClient, conf, accessToken, refreshToken, observers...)
		if err != nil {
			return err
		}
		r.httpClient = httpClient
		return nil
	}
}

// WithClient sets the http client to use for requests.
func WithAccountUUID(accountUUID string) Option {
	return func(r *Requester) error {
		if r.headers == nil {
			r.headers = make(http.Header)
		}
		r.headers[meroxaAccountUUIDHeader] = []string{accountUUID}
		return nil
	}
}
