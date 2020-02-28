package cmd

import "fmt"

type Config map[string]string

func (c Config) Set(key, value string) error {
	c[key] = value
	return nil
}

func (c Config) Get(key string) (string, bool) {
	v, ok := c[key]
	return v, ok
}

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
