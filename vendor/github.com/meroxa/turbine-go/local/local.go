package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"time"
	"unsafe"

	"github.com/meroxa/turbine-go"
)

type Turbine struct {
	config turbine.AppConfig
}

func New() Turbine {
	ac, err := turbine.ReadAppConfig("", "")
	if err != nil {
		log.Fatalln(err)
	}
	return Turbine{ac}
}

func (t Turbine) Resources(name string) (turbine.Resource, error) {
	return Resource{
		Name:         name,
		fixturesPath: t.config.Resources[name],
	}, nil
}

func (t Turbine) Process(rr turbine.Records, fn turbine.Function) turbine.Records {
	var out turbine.Records

	// use reflection to access intentionally hidden fields
	inVal := reflect.ValueOf(&rr).Elem().FieldByName("records")

	// hack to create reference that can be accessed
	in := reflect.NewAt(inVal.Type(), unsafe.Pointer(inVal.UnsafeAddr())).Elem()
	inRR := in.Interface().([]turbine.Record)

	rawOut := fn.Process(inRR)
	out = turbine.NewRecords(rawOut)

	return out
}

type Resource struct {
	Name         string
	fixturesPath string
}

func (r Resource) Records(collection string, cfg turbine.ResourceConfigs) (turbine.Records, error) {
	binPath, err := os.Executable()
	if err != nil {
		return turbine.Records{}, err
	}
	if r.fixturesPath == "" {
		return turbine.Records{},
			fmt.Errorf("must specify fixtures path to data for source resources in order to run locally")
	}
	dirPath := path.Dir(binPath)
	pwd := fmt.Sprintf("%s/%s", dirPath, r.fixturesPath)
	return readFixtures(pwd, collection)
}

func (r Resource) WriteWithConfig(rr turbine.Records, collection string, cfg turbine.ResourceConfigs) error {
	prettyPrintRecords(r.Name, collection, turbine.GetRecords(rr))
	return nil
}

func (r Resource) Write(rr turbine.Records, collection string) error {
	return r.WriteWithConfig(rr, collection, turbine.ResourceConfigs{})
}

func prettyPrintRecords(name string, collection string, rr []turbine.Record) {
	fmt.Printf("=====================to %s (%s) resource=====================\n", name, collection)
	for _, r := range rr {
		payloadVal := string(r.Payload)
		m, err := r.Payload.Map()
		if err == nil {
			b, err := json.MarshalIndent(m, "", "    ")
			if err == nil {
				payloadVal = string(b)
			}
		}
		fmt.Println(payloadVal)
	}
	fmt.Printf("%d record(s) written\n", len(rr))
}

type fixtureRecord struct {
	Key       string
	Value     map[string]interface{}
	Timestamp string
}

func readFixtures(path, collection string) (turbine.Records, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return turbine.Records{}, err
	}

	var records map[string][]fixtureRecord
	err = json.Unmarshal(b, &records)
	if err != nil {
		return turbine.Records{}, err
	}

	var rr []turbine.Record
	for _, r := range records[collection] {
		rr = append(rr, wrapRecord(r))
	}

	return turbine.NewRecords(rr), nil
}

func wrapRecord(m fixtureRecord) turbine.Record {
	b, _ := json.Marshal(m.Value)

	var t time.Time
	if m.Timestamp == "" {
		t = time.Now()
	} else {
		t, _ = time.Parse(time.RFC3339, m.Timestamp)
	}

	return turbine.Record{
		Key:       m.Key,
		Payload:   b,
		Timestamp: t,
	}
}

// RegisterSecret pulls environment variables with the same name
func (t Turbine) RegisterSecret(name string) error {
	val := os.Getenv(name)
	if val == "" {
		return errors.New("secret is invalid or not set")
	}

	return nil
}
