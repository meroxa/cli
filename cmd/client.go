package cmd

import (
	"github.com/meroxa/meroxa-go"
)

func client() (*meroxa.Client, error) {
	u, p, err := readCreds()
	if err != nil {
		return nil, err
	}
	return meroxa.New(u, p, versionString())
}
