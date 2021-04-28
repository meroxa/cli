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

		// create subdirectory if it doesn't exist, otherwise viper will complain
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("could not create meroxa config directory: %w", err)
		}

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

		// No config found, fallback to deprecated config file location in $HOME
		// TODO remove this code once we migrate acceptance tests to use new location
		if err := setupCompatibility(cfg); err != nil {
			return nil, err
		}
	}

	// TODO remove this code once we migrate acceptance tests to use new env variable
	if apiURL, ok := os.LookupEnv("API_URL"); ok {
		os.Setenv("MEROXA_API_URL", apiURL)
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

// setupCompatibility falls back to deprecated config file location in $HOME
// also it enables env variable API_URL alongside MEROXA_API_URL.
// This function should be removed once we migrate acceptance tests.
func setupCompatibility(cfg *viper.Viper) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %w", err)
	}
	cfg.AddConfigPath(homeDir)
	cfg.SetConfigName("meroxa")
	cfg.SetConfigType("env")

	if err = cfg.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("could not read config: %w", err)
		}
	}
	cfg.SetConfigName("config") // revert config name
	if err == nil {
		// we read the config in the home folder, let's write it to the new location
		err = cfg.SafeWriteConfig()
		if err != nil {
			return err
		}
	}

	return nil
}