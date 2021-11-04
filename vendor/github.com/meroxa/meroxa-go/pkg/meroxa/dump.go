package meroxa

import (
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

type dumpTransport struct {
	out                    io.Writer
	transport              http.RoundTripper
	obfuscateAuthorization bool
}

func (d *dumpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := d.dumpRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := d.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	err = d.dumpResponse(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *dumpTransport) dumpRequest(req *http.Request) error {
	if d.obfuscateAuthorization {
		if auth := req.Header.Get("Authorization"); auth != "" {
			save := auth
			defer func() {
				// restore old header
				req.Header.Set("Authorization", save)
			}()

			tokens := strings.SplitN(auth, " ", 2)
			if len(tokens) == 2 {
				tokens[1] = d.obfuscate(tokens[1])
				auth = strings.Join(tokens, " ")
			}
			req.Header.Set("Authorization", auth)
		}
	}

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}

	_, err = d.out.Write(dump)
	if err != nil {
		return err
	}
	return nil
}

func (d *dumpTransport) dumpResponse(resp *http.Response) error {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	_, err = d.out.Write(dump)
	if err != nil {
		return err
	}
	return nil
}

func (d *dumpTransport) obfuscate(text string) string {
	if len(text) < 5 {
		// hide whole text
		return strings.Repeat("*", len(text))
	}

	const (
		maxVisibleLen = 7
	)

	visibleLen := (len(text) - 3) / 2
	if visibleLen > maxVisibleLen {
		visibleLen = maxVisibleLen
	}

	return text[:visibleLen] + "..." + text[len(text)-visibleLen:]
}
