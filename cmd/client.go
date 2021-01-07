package cmd

import (
	"github.com/meroxa/cli/credstore"
	"github.com/meroxa/meroxa-go"
)

func client() (*meroxa.Client, error) {
	_, token, err := credstore.Get(meroxaLabel, meroxaURL)
	if err != nil {
		return nil, err
	}
	return meroxa.New(token, versionString())
}
