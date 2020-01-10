package cmd

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

type meroxaAPIClient struct {
	APIEndpoint *url.URL
	*http.Client
}

func newClient(urlString string) *meroxaAPIClient {
	apiEndpoint, err := meroxaAPIURL(urlString)
	if err != nil {
		log.Fatal("invalid Meroxa API URL provided")
	}

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &meroxaAPIClient{
		APIEndpoint: apiEndpoint,
		Client:      c,
	}
}

func meroxaAPIURL(urlString string) (*url.URL, error) {
	if urlString == "" {
		urlString = "https://api.meroxa.io/v1/"
	}

	return url.Parse(urlString)
}
