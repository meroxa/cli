/*
Copyright © 2022 Meroxa Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package global

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	// The environment variable prefix of all environment variables bound to our command line flags.
	envPrefix = "MEROXA"
	envName   = "config"
	envType   = "env"
)

func GetMeroxaAPIURL() string {
	return getEnvVal([]string{"MEROXA_API_URL"}, "https://api.meroxa.io")
}

func GetMeroxaAuthAudience() string {
	return getEnvVal([]string{"MEROXA_AUTH_AUDIENCE", "MEROXA_AUDIENCE"}, "https://api.meroxa.io/v1")
}

func GetMeroxaAuthDomain() string {
	return getEnvVal([]string{"MEROXA_AUTH_DOMAIN", "MEROXA_DOMAIN"}, "auth.meroxa.io")
}

func GetMeroxaAuthClientID() string {
	return getEnvVal([]string{"MEROXA_AUTH_CLIENT_ID", "MEROXA_CLIENT_ID"}, "2VC9z0ZxtzTcQLDNygeEELV3lYFRZwpb")
}

func GetLocalTurbineJSSetting() string {
	return getEnvVal([]string{"MEROXA_USE_LOCAL_TURBINE_JS"}, "false")
}

func getMeroxaAuthCallbackPort() string {
	return getEnvVal([]string{MeroxaAuthCallbackPort}, "21900")
}

// getMeroxaAuthCallbackHost will return the callback host.
// Note: If port is desired, it'll need to be included here.
func getMeroxaAuthCallbackHost() string {
	defaultHost := fmt.Sprintf("localhost:%s", getMeroxaAuthCallbackPort())
	// check if port was included or not
	return getEnvVal([]string{MeroxaAuthCallbackHost}, defaultHost)
}

func getMeroxaAuthCallbackProtocol() string {
	return getEnvVal([]string{MeroxaAuthCallbackProtocol}, "http")
}

// GetMeroxaAuthCallbackURL will return either the user configured oauth callback url
// or a default one: "http://localhost:21900/oauth/callback".
func GetMeroxaAuthCallbackURL() string {
	callback := url.URL{
		Scheme: getMeroxaAuthCallbackProtocol(),
		Host:   getMeroxaAuthCallbackHost(),
		Path:   "/oauth/callback",
	}
	return getEnvVal([]string{MeroxaAuthCallbackURL}, callback.String())
}

// getEnvVal returns the value of either the first existing key specified in keys, or defaultVal if none were present.
func getEnvVal(keys []string, defaultVal string) string {
	for _, key := range keys {
		if Config != nil {
			// First tries to return the value from the meroxa configuration file
			if val := Config.GetString(key); val != "" {
				return val
			}
		}

		/// Tries to fetch it from the environment if not.
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
	}
	return defaultVal
}

func readConfig() (*viper.Viper, error) {
	cfg := viper.New()

	if flagCLIConfigFile != "" {
		// Use config file from the flag.
		cfg.SetConfigFile(flagCLIConfigFile)
	} else {
		// Find home directory.
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("could not get config directory: %w", err)
		}
		configDir = filepath.Join(configDir, "meroxa")

		// create subdirectory if it doesn't exist, otherwise viper will complain
		err = os.MkdirAll(configDir, 0o755)
		if err != nil {
			return nil, fmt.Errorf("could not create meroxa config directory: %w", err)
		}

		cfg.AddConfigPath(configDir)
		cfg.SetConfigName(envName)
		cfg.SetConfigType(envType)
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := cfg.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("could not read config: %w", err)
		}

		// No config found, fallback to old config file location in $HOME
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

// setupCompatibility falls back to old config file location in $HOME
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
