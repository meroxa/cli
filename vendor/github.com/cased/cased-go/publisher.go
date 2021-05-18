package cased

import (
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var currentPublisher Publisher = NewPublisher()

// PublisherOptions ...
type PublisherOptions struct {
	// PublishURL contains the URL to published Cased events to.
	PublishURL string `envconfig:"CASED_PUBLISH_URL" default:"https://publish.cased.com"`

	// PublishKey is the publish API key used to publish to an audit trail.
	//
	// A publish key is associated with a single audit trail and is required if
	// you intend to publish events to Cased in your application.
	PublishKey string `envconfig:"CASED_PUBLISH_KEY"`

	Debug bool `envconfig:"CASED_DEBUG" default:"false"`

	// Silence to determine if new events are published to Cased.
	Silence bool `envconfig:"CASED_SILENCE" default:"false"`

	HTTPClient    *http.Client
	HTTPTransport *http.Transport
	HTTPTimeout   time.Duration `envconfig:"CASED_HTTP_TIMEOUT" default:"5s"`

	Transport Transporter
}

// PublisherOption ...
type PublisherOption func(opts *PublisherOptions)

// Client is the underlying processor that is used by the main API and Hub
// instances. It must be created with NewClient.
type Client struct {
	options   PublisherOptions
	transport Transporter
}

// CurrentPublisher ...
func CurrentPublisher() Publisher {
	return currentPublisher
}

// SetPublisher ...
func SetPublisher(publisher Publisher) {
	currentPublisher = publisher
}

// NewPublisher ...
func NewPublisher(opts ...PublisherOption) Publisher {
	publisherOpts := PublisherOptions{}
	_ = envconfig.Process("", &publisherOpts)

	for _, opt := range opts {
		opt(&publisherOpts)
	}

	if publisherOpts.Debug {
		Logger.SetOutput(os.Stderr)
	}

	client := &Client{
		options: publisherOpts,
	}

	client.setupTransport()

	return client
}

// WithPublishKey configures the publish key used to publish audit events to
// Cased. You can obtain the publish key from the audit trail settings page
// within the Cased dashboard.
func WithPublishKey(publishKey string) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.PublishKey = publishKey
	}
}

// WithPublishURL ...
func WithPublishURL(publishURL string) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.PublishURL = publishURL
	}
}

// WithSilence ...
func WithSilence(silence bool) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.Silence = silence
	}
}

// WithHTTPClient ...
func WithHTTPClient(httpClient *http.Client) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.HTTPClient = httpClient
	}
}

// WithHTTPTransport ...
func WithHTTPTransport(httpTransport *http.Transport) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.HTTPTransport = httpTransport
	}
}

// WithHTTPTimeout ...
func WithHTTPTimeout(httpTimeout time.Duration) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.HTTPTimeout = httpTimeout
	}
}

// WithTransport ...
func WithTransport(transport Transporter) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.Transport = transport
	}
}

// WithDebug ...
func WithDebug(debug bool) PublisherOption {
	return func(opts *PublisherOptions) {
		opts.Debug = debug
	}
}

// Options return PublisherOptions for the current Client.
func (c Client) Options() PublisherOptions {
	return c.options
}

// Publish ...
func (c Client) Publish(event AuditEvent) error {
	aep := NewAuditEventPayload(event)

	return c.transport.Publish(aep)
}

// Flush ...
func (c *Client) Flush(timeout time.Duration) bool {
	return c.transport.Flush(timeout)
}

func (c *Client) setupTransport() {
	opts := c.options
	transport := opts.Transport

	if transport == nil {
		if opts.PublishKey == "" {
			Logger.Print("No publish key detected, no audit events will be published to Cased. Set CASED_PUBLISH_KEY to publish audit events.")
			transport = NewNoopHTTPTransport()
		} else {
			transport = NewHTTPTransport()
		}
	}

	transport.Configure(opts)
	c.transport = transport
}
