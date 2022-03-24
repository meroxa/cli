package turbine

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Records struct {
	Stream  string
	records []Record
}

func NewRecords(rr []Record) Records {
	return Records{records: rr}
}

func GetRecords(r Records) []Record {
	return r.records
}

type RecordsWithErrors struct {
	Stream  string
	records []RecordWithError
}

type Record struct {
	Key       string
	Payload   Payload
	Timestamp time.Time
}

// JSONSchema returns true if the record is formatted with JSON Schema, false otherwise
func (r Record) JSONSchema() bool {
	p, err := r.Payload.Map()
	if err != nil {
		return false
	}

	if _, ok := p["schema"]; ok {
		if _, ok := p["payload"]; ok {
			return true
		}
		return false
	}

	return false
}

type Payload []byte

func (p Payload) Map() (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(p, &m)
	return m, err
}

func (p Payload) Get(path string) interface{} {
	nestedPath := strings.Join([]string{"payload", path}, ".")
	return gjson.Get(string(p), nestedPath).Value()
}

// TODO: Add GetType(path string) to tell you what the data type is.
// TODO: Should we passthrough the gjson helper methods?

type schemaField struct {
	Field   string `json:"field"`
	Options bool   `json:"optional"`
	Type    string `json:"type"`
}

func (p *Payload) Set(path string, value interface{}) error {
	// update payload
	nestedPath := strings.Join([]string{"payload", path}, ".")
	val, err := sjson.Set(string(*p), nestedPath, value)
	if err != nil {
		return err
	}
	*p = []byte(val)

	// update schema
	schemaField := map[string]string{
		"field":    path,
		"optional": "true",
		"type":     "string", // TODO: map Go types to JSON types
	}

	schemaNestedPath := strings.Join([]string{"schema", "fields.-1"}, ".")
	sval, err := sjson.Set(string(*p), schemaNestedPath, schemaField)
	if err != nil {
		return err
	}
	*p = []byte(sval)

	return nil
}

type RecordWithError struct {
	Error error
	Record
}
