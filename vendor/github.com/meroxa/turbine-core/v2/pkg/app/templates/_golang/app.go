package main

import (
	// Dependencies of the example data app
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	// Dependencies of Turbine
	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/meroxa/turbine-go/v3/pkg/turbine"
	"github.com/meroxa/turbine-go/v3/pkg/turbine/cmd"
)

func main() {
	cmd.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(t turbine.Turbine) error {
	// Identify an upstream data store for your data app
	// with the `Source` function
	// Replace `source_name` with the source name of your choice
	// and use the desired Conduit Plugin you'd like to use.

	pg, err := t.Source("source_name", "postgres", turbine.WithPluginConfig(map[string]string{
		"url":                     "url",
		"key":                     "key",
		"table":                   "records",
		"columns":                 "key,column1,column2,column3",
		"cdcMode":                 "logrepl",
		"logrepl.publicationName": "meroxademo",
		"logrepl.slotName":        "meroxademo",
	}))
	if err != nil {
		return err
	}

	records, err := pg.Read()
	if err != nil {
		return err
	}

	// Specify what code to execute against upstream records
	// with the `Process` function
	// Replace `Anonymize` with the name of your function code.

	processed, err := t.Process(records, Anonymize{})
	if err != nil {
		return err
	}

	// Identify a downstream data store for your data app
	// with the `Destination` function
	// Replace `destination_name` with the destination name of your choice
	// and use the desired Conduit Plugin you'd like to use.

	s3, err := t.Destination("destination_name", "s3",
		turbine.WithPluginConfig(
			map[string]string{
				"aws.accessKeyId":     "id",
				"aws.secretAccessKey": "key",
				"aws.region":          "us-east-1",
				"aws.bucket":          "bucket_name",
			}))
	if err != nil {
		return err
	}
	err = s3.Write(processed)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(stream []opencdc.Record) []opencdc.Record {
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
