package jsonpath

import (
	"fmt"
	"reflect"
	"sync"
	"unicode"
)

// Delimiter ...
const Delimiter = "."

// Reader ...
type Reader struct {
	json      map[string]interface{}
	table     map[string]interface{}
	paths     []string
	processed bool
	mu        sync.Mutex
}

// NewReader ...
func NewReader(json map[string]interface{}) *Reader {
	return &Reader{
		json:      json,
		table:     map[string]interface{}{},
		paths:     []string{},
		processed: false,
	}
}

// Paths ...
func (r *Reader) Paths() []string {
	r.generate()

	return r.paths
}

// Path ...
func (r *Reader) Path(path string) interface{} {
	r.generate()

	return r.table[path]
}

func (r *Reader) generate() {
	if r.processed {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	v := reflect.ValueOf(r.json)
	r.walk(v, "")

	paths := make([]string, 0, len(r.table))
	for path := range r.table {
		paths = append(paths, path)
	}
	r.paths = paths

	r.processed = true
}

func (r *Reader) walk(v reflect.Value, path string) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			p := fmt.Sprintf("%s[%d]", path, i)
			r.walk(v.Index(i), p)
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			key := escapeKey(iter.Key().String())
			nestedPath := key
			if path != "" {
				nestedPath = fmt.Sprintf("%s%s%s", path, Delimiter, key)
			}
			r.walk(iter.Value(), nestedPath)
		}
	case reflect.Invalid:
		// nil value
	default:
		r.buildPath(path, v.Interface())
	}
}

func (r *Reader) buildPath(path string, value interface{}) {
	// Key is the path to the JSON property we're indexing.
	v := fmt.Sprintf("%s%s", Delimiter, path)

	// We've already seen this path and value, skip...
	if _, ok := r.table[v]; ok {
		return
	}

	r.table[v] = value
}

func escapeKey(key string) string {
	escape := false
	for _, v := range key {
		// Test each character to see if it is whitespace.
		if unicode.IsSpace(v) {
			escape = true
			break
		}
	}

	if escape {
		return fmt.Sprintf("\"%s\"", key)
	}

	return key
}
