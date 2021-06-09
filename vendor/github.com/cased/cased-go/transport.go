package cased

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const defaultBufferSize = 100
const defaultTimeout = 10 * time.Second

// Transporter ...
type Transporter interface {
	Configure(options PublisherOptions)
	Publish(event *AuditEventPayload) error
	Flush(timeout time.Duration) bool
}

type batch struct {
	events  chan *AuditEventPayload
	started chan struct{}
	done    chan struct{}
}

// HTTPTransport ...
type HTTPTransport struct {
	client    *http.Client
	transport *http.Transport
	timeout   time.Duration

	BufferSize int

	buffer chan batch

	start sync.Once
}

// NewHTTPTransport ...
func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{
		BufferSize: defaultBufferSize,
	}
}

// Configure prepares the asynchronous audit event publisher with provided
// client options.
func (t *HTTPTransport) Configure(options PublisherOptions) {
	// Heavily influenced by Sentry's flush pattern, requirements align closely.
	t.buffer = make(chan batch, 1)

	// Prepare buffer with its first batch, used to obtain batch when publishing
	// first audit event.
	t.buffer <- batch{
		events:  make(chan *AuditEventPayload, t.BufferSize),
		started: make(chan struct{}),
		done:    make(chan struct{}),
	}

	if options.HTTPTransport != nil {
		t.transport = options.HTTPTransport
	} else {
		t.transport = &http.Transport{}
	}

	if options.HTTPTimeout > 0 {
		t.timeout = options.HTTPTimeout
	} else {
		t.timeout = defaultTimeout
	}

	if options.HTTPClient != nil {
		t.client = options.HTTPClient
	} else {
		t.client = &http.Client{
			Transport: t.transport,
			Timeout:   t.timeout,
		}
	}

	t.start.Do(func() {
		go t.worker()
	})
}

// Flush waits for all audit events to be published that are in the buffer.
func (t *HTTPTransport) Flush(timeout time.Duration) bool {
	expired := time.After(timeout)

	for {
		select {
		case b := <-t.buffer:
			select {
			case <-b.started:
				close(b.events)

				t.buffer <- batch{
					events:  make(chan *AuditEventPayload, t.BufferSize),
					started: make(chan struct{}),
					done:    make(chan struct{}),
				}

				select {
				case <-b.done:
					Logger.Println("Published all audit events in buffer.")
					return true
				case <-expired:
					Logger.Printf("Could not flush all audit events from buffer, timed out after %s.\n", timeout.String())
					return false
				}

			default:
				// Put buffer back until it has started
				t.buffer <- b
			}
		case <-expired:
			Logger.Printf("Could not flush all audit events from buffer, timed out after %s.\n", timeout.String())
			return false
		}
	}
}

// Publish queues the audit event to be published in the asynchronously.
//
// To ensure queued audit events are published at end of process see Flush.
func (t *HTTPTransport) Publish(event *AuditEventPayload) error {
	// Obtain the buffer lock
	b := <-t.buffer

	// Add the event to the buffer
	b.events <- event

	// Release buffer lock
	t.buffer <- b

	return nil
}

func (t *HTTPTransport) worker() {
	for b := range t.buffer {
		// Signal batch has started processing.
		close(b.started)

		// Release lock on buffer.
		t.buffer <- b

		// Publish all audit events to Cased based on client's configuration.
		for event := range b.events {
			_, err := publish(t.client, event)
			if err != nil {
				Logger.Printf("There was an issue with publishing audit event: %v", err)
				continue
			}
		}

		// Signal that processing of the batch is done. Useful for when flushing
		// audit events at end of process.
		close(b.done)
	}
}

// HTTPSyncTransport provides a transport that publishes audit events
// synchronously as they are received.
type HTTPSyncTransport struct {
	client    *http.Client
	transport *http.Transport
	timeout   time.Duration
}

// NewHTTPSyncTransport returns a transport that publishes audit events
// synchronously as they are received.
func NewHTTPSyncTransport() *HTTPSyncTransport {
	return &HTTPSyncTransport{}
}

// Configure prepares the synchronous audit event publisher with provided client
// options.
func (t *HTTPSyncTransport) Configure(options PublisherOptions) {
	if options.HTTPTransport != nil {
		t.transport = options.HTTPTransport
	} else {
		t.transport = &http.Transport{}
	}

	if options.HTTPTimeout > 0 {
		t.timeout = options.HTTPTimeout
	} else {
		t.timeout = time.Second * 30
	}

	if options.HTTPClient != nil {
		t.client = options.HTTPClient
	} else {
		t.client = &http.Client{
			Transport: t.transport,
			Timeout:   t.timeout,
		}
	}
}

// Flush is unused.
func (t *HTTPSyncTransport) Flush(_ time.Duration) bool {
	return true
}

// Publish publishes the provided audit event to Cased.
func (t *HTTPSyncTransport) Publish(event *AuditEventPayload) error {
	_, err := publish(t.client, event)
	return err
}

// NoopHTTPTransport does not publish audit events to Cased.
type NoopHTTPTransport struct{}

// NewNoopHTTPTransport returns a client that does not publish audit events to
// Cased.
func NewNoopHTTPTransport() *NoopHTTPTransport {
	return &NoopHTTPTransport{}
}

// Configure is a noop operation.
func (t *NoopHTTPTransport) Configure(options PublisherOptions) {
}

// Publish is a noop operation.
func (t *NoopHTTPTransport) Publish(event *AuditEventPayload) error {
	return nil
}

// Flush is a noop operation.
func (t *NoopHTTPTransport) Flush(_ time.Duration) bool {
	return true
}

func publish(client *http.Client, event *AuditEventPayload) (*http.Response, error) {
	p := CurrentPublisher()
	body, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, p.Options().PublishURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// TODO: Use build flags to encode git sha in user agent
	req.Header.Set("User-Agent", "cased-go/v0.1")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.Options().PublishKey))

	resp, err := client.Do(req)
	if err != nil {
		Logger.Print("Could not publish event")
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return resp, errors.New("unauthorized")
	case http.StatusOK, http.StatusCreated:
		Logger.Println("Successfully published audit event.")
		return resp, nil
	default:
		return resp, fmt.Errorf("Received %s while publishing audit event", resp.Status)
	}
}
