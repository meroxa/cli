package cmd

import "fmt"

// Config defines a dictionary to be used across our cmd package
type Config map[string]string

// Set a key value pair
func (c Config) Set(key, value string) error {
	c[key] = value
	return nil
}

// Get a key value pair
func (c Config) Get(key string) (string, bool) {
	v, ok := c[key]
	return v, ok
}

// Merge one Config definition onto another
func (c Config) Merge(cfg Config) error {
	for k, v := range cfg {
		_, exist := c[k]
		if exist {
			return fmt.Errorf("merge config, key %s already present", k)
		}
		c[k] = v
	}
	return nil
}
