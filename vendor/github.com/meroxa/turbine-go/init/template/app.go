package main

import (
// Dependencies of the example data app
	"crypto/md5"
	"encoding/hex"
	"log"
	
// Dependencies of Turbine
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	// To configure your data stores as resources on the Meroxa Platform
	// use the Meroxa Dashboard, CLI, or Meroxa Terraform Provider
	// For more details refer to: https://docs.meroxa.com/

	// Identify an upstream data store for your data app
	// with the `Resources` function
	// Replace `source_name` with the resource name the
	// data store was configured with on Meroxa
	source, err := v.Resources("source_name")
	if err != nil {
		return err
	}

	// Specify which upstream records to pull
	// with the `Records` function
	// Replace `collection_name` with a table, collection,
	// or bucket name in your data store
	rr, err := source.Records("collection_name", nil)
	if err != nil {
		return err
	}

	// Specify what code to execute against upstream records
	// with the `Process` function
	// Replace `Anonymize` with the name of your function code
	res, _ := v.Process(rr, Anonymize{})

	// Identify a downstream data store for your data app
	// with the `Resources` function
	// Replace `dest_name` with the resource name the
	// data store was configured with on Meroxa
	dest, err := v.Resources("dest_name")
	if err != nil {
		return err
	}

	// Specify where to write records downstream
	// using the `Write` function
	// Replace `collection_name` with a table, collection,
	// or bucket name in your data store
	err = dest.Write(res, "collection_name", nil)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(stream []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	for i, r := range stream {
		hashedEmail := consistentHash(r.Payload.Get("email").(string))
		err := r.Payload.Set("email", hashedEmail)
		if err != nil {
			log.Println("error setting value: ", err)
			break
		}
		stream[i] = r
	}
	return stream, nil
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}