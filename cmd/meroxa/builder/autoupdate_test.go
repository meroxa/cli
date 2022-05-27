package builder

import (
	"testing"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/spf13/viper"
)

func TestNeedToCheckNewerCLIVersion(t *testing.T) {
	oldConfig := global.Config

	tests := []struct {
		desc     string
		config   func() *viper.Viper
		expected bool
	}{
		{
			desc: "Global config is nil",
			config: func() *viper.Viper {
				return nil
			},
			expected: false,
		},
		{
			desc: "Notifications are disabled by the user",
			config: func() *viper.Viper {
				cfg := viper.New()
				cfg.Set(global.DisableNotificationsUpdate, "true")
				return cfg
			},
			expected: false,
		},
		{
			desc:     "Version has never been checked",
			config:   viper.New,
			expected: true,
		},
		{
			desc: "Version was recently checked",
			config: func() *viper.Viper {
				cfg := viper.New()
				// Checked last time less than a week before
				lastTimeChecked := time.Now().UTC().AddDate(0, 0, -7).Add(time.Minute * 1)
				cfg.Set(global.LatestCLIVersionUpdatedAtEnv, lastTimeChecked)
				return cfg
			},
			expected: false,
		},
		{
			desc: "Version was checked more than a week before",
			config: func() *viper.Viper {
				cfg := viper.New()
				// Checked last time more than a week before
				lastTimeChecked := time.Now().UTC().AddDate(0, 0, -7).Add(-time.Minute * 1)
				cfg.Set(global.LatestCLIVersionUpdatedAtEnv, lastTimeChecked)
				return cfg
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			global.Config = tc.config()
			got := needToCheckNewerCLIVersion()

			if tc.expected != got {
				t.Fatalf("expected needToCheckNewerCLIVersion to be %v, got %v", tc.expected, got)
			}
		})
	}

	global.Config = oldConfig
}
