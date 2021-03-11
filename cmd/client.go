package cmd

import (
	"os"

	"github.com/meroxa/meroxa-go"
)

func isDebugEnabled() bool {
	if val, ok := os.LookupEnv("MEROXA_DEBUG"); ok {
		if val == "1" {
			return true
		}
	}

	return false
}

func client() (*meroxa.Client, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}

	return meroxa.New(accessToken, versionString(), isDebugEnabled())
}
