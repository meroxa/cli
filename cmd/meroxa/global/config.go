package global

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	// The name of our config file, without the file extension because viper supports many different config file languages.
	defaultConfigFilename = "meroxa"

	// The environment variable prefix of all environment variables bound to our command line flags.
	envPrefix = "MEROXA"
)

func readConfig() (*viper.Viper, error) {
	cfg := viper.New()

	if flagConfig != "" {
		// Use config file from the flag.
		cfg.SetConfigFile(flagConfig)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return nil, fmt.Errorf("could not get home directory: %w", err)
		}

		// Set the base name of the config file, without the file extension.
		cfg.SetConfigName(defaultConfigFilename)
		cfg.AddConfigPath(home)
	}
	cfg.SetConfigType("env")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := cfg.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("could not read config: %w", err)
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable MEROXA_NUMBER. This helps
	// avoid conflicts.
	cfg.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	cfg.AutomaticEnv()

	return cfg, nil
}
