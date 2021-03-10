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
	dump, _ := httputil.DumpRequestOut(h, true)
	log.Println(string(dump))

	resp, err := d.r.RoundTrip(h)
	dump, _ = httputil.DumpResponse(resp, true)
	log.Println(string(dump))

	log.Println(resp)
	return resp, err
}

func httpDebugClient() *http.Client {
	return &http.Client{
		Transport: &DumpTransport{http.DefaultTransport},
		Timeout:   ClientTimeOut,
	}
}
