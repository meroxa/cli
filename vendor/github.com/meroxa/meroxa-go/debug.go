package meroxa

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type DumpTransport struct {
	r http.RoundTripper
}

func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	cloned := h.Clone(context.Background())

	// Makes sure we don't log out the bearer token by accident when it's not nil
	if !strings.Contains(cloned.Header.Get("Authorization"), "nil") {
		cloned.Header.Set("Authorization", "REDACTED")
	}

	dump, err := httputil.DumpRequestOut(cloned, true)

	if err != nil {
		return nil, err
	}

	log.Printf(string(dump))

	resp, err := d.r.RoundTrip(h)

	if err != nil {
		return nil, err
	}

	dump, err = httputil.DumpResponse(resp, true)

	if err != nil {
		return nil, err
	}

	log.Printf(string(dump))

	return resp, err
}
