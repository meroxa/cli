package meroxa

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type DumpTransport struct {
	r http.RoundTripper
}

func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	dump, err := httputil.DumpRequestOut(h, true)

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
