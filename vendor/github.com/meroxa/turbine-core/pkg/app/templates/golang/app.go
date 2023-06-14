package main

import (
	// Dependencies of the example data app
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	// Dependencies of Turbine
	"github.com/meroxa/turbine-go/v2/pkg/turbine"
	"github.com/meroxa/turbine-go/v2/pkg/turbine/cmd"
)

func main() {
	cmd.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	// To configure your data stores as resources on the Meroxa Platform
	// use the Meroxa Dashboard, CLI, or Meroxa Terraform Provider.
	// For more details refer to: https://docs.meroxa.com/
	//
	// Identify an upstream data store for your data app
	// with the `Resources` function
	// Replace `source_name` with the resource name the
	// data store was configured with on Meroxa.

	source, err := v.Resources("source_name")
	if err != nil {
		return err
	}

	// Specify which upstream records to pull
	// with the `Records` function
	// Replace `collection_name` with a table, collection,
	// or bucket name in your data store.
	// If a configuration is needed for your source,
	// you can pass it as a second argument to the `Records` function. For example:
	//
	// source.Records("collection_name", turbine.ConnectionOptions{turbine.ResourceConfig{Field: "incrementing.field.name", Value:"id"}})

	rr, err := source.Records("collection_name", nil)
	if err != nil {
		return err
	}

	// Specify what code to execute against upstream records
	// with the `Process` function
	// Replace `Anonymize` with the name of your function code.

	res, err := v.Process(rr, Anonymize{})
	if err != nil {
		return err
	}

	// Identify a downstream data store for your data app
	// with the `Resources` function
	// Replace `destination_name` with the resource name the
	// data store was configured with on Meroxa.

	dest, err := v.Resources("destination_name")
	if err != nil {
		return err
	}

	// Specify where to write records downstream
	// using the `Write` function
	// Replace `collection_archive` with a table, collection,
	// or bucket name in your data store.
	// If a configuration is needed, you can also use i.e.
	//
	// dest.WriteWithConfig(
	//  res,
	//  "my-archive",
	//  turbine.ConnectionOptions{turbine.ResourceConfig{Field: "buffer.flush.time", Value: "10"}}
	// )

	err = dest.Write(res, "collection_archive")
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(stream []turbine.Record) []turbine.Record {
	for i, record := range stream {
		email := fmt.Sprintf("%s", record.Payload.Get("after.customer_email"))
		if email == "" {
			log.Printf("unable to find customer_email value in record %d\n", i)
			break
		}
		hashedEmail := consistentHash(email)
		err := record.Payload.Set("after.customer_email", hashedEmail)
		if err != nil {
			log.Println("error setting value: ", err)
			continue
		}
		stream[i] = record
	}
	return stream
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
