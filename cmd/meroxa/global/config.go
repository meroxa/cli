package global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
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
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("could not get config directory: %w", err)
		}
		configDir = filepath.Join(configDir, "meroxa")

		cfg.AddConfigPath(configDir)
		cfg.SetConfigName("config")
		cfg.SetConfigType("env")
	}

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

	// Add support for flags like --favorite-color by replacing - with _.
	cfg.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables.
	cfg.AutomaticEnv()

	return cfg, nil
}
