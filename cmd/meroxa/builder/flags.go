package builder

import (
	"fmt"
	"reflect"
	"strconv"
)

func BuildFlags(obj interface{}) []Flag {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		panic(fmt.Errorf("expected a pointer, got %s", v.Kind()))
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected a struct, got %s", v.Kind()))
	}
	t := v.Type()

	var err error
	flags := make([]Flag, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		flags[i], err = buildFlag(v.Field(i), t.Field(i))
		if err != nil {
			panic(err)
		}
	}
	return flags
}

func buildFlag(val reflect.Value, sf reflect.StructField) (Flag, error) {
	const (
		tagNameLong       = "long"
		tagNameShort      = "short"
		tagNameRequired   = "required"
		tagNamePersistent = "persistent"
		tagNameUsage      = "usage"
		tagNameHidden     = "hidden"
		tagNameDefault    = "default"
	)

	var (
		long       string
		short      string
		required   bool
		persistent bool
		usage      string
		hidden     bool
		def        string
	)

	if v, ok := sf.Tag.Lookup(tagNameLong); ok {
		long = v
	}
	if v, ok := sf.Tag.Lookup(tagNameShort); ok {
		short = v
	}
	if v, ok := sf.Tag.Lookup(tagNameRequired); ok {
		var err error
		required, err = strconv.ParseBool(v)
		if err != nil {
			return Flag{}, fmt.Errorf("error parsing tag \"required\": %w", err)
		}
	}
	if v, ok := sf.Tag.Lookup(tagNamePersistent); ok {
		var err error
		persistent, err = strconv.ParseBool(v)
		if err != nil {
			return Flag{}, fmt.Errorf("error parsing tag \"persistent\": %w", err)
		}
	}
	if v, ok := sf.Tag.Lookup(tagNameUsage); ok {
		usage = v
	}
	if v, ok := sf.Tag.Lookup(tagNameHidden); ok {
		var err error
		hidden, err = strconv.ParseBool(v)
		if err != nil {
			return Flag{}, fmt.Errorf("error parsing tag \"hidden\": %w", err)
		}
	}
	if v, ok := sf.Tag.Lookup(tagNameDefault); ok {
		var err error
		def = v
		if err != nil {
			return Flag{}, fmt.Errorf("error parsing tag \"default\": %w", err)
		}
	}

	return Flag{
		Long:       long,
		Short:      short,
		Usage:      usage,
		Required:   required,
		Persistent: persistent,
		Default:    def,
		Ptr:        val.Addr().Interface(),
		Hidden:     hidden,
	}, nil
}
