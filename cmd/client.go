package cmd

import (
	"github.com/meroxa/meroxa-go"
)

func client() (*meroxa.Client, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}

	return meroxa.New(accessToken, versionString())
}
