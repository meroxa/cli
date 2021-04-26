/*
Copyright © 2020 Meroxa Inc

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

package old

import (
	"fmt"
	"os"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const (
	DashboardProductionURL = "https://dashboard.meroxa.io"
	DashboardStagingURL    = "https://dashboard.staging.meroxa.io"
)

func getBillingURL() string {
	platformURL := DashboardProductionURL

	if os.Getenv("ENV") == "staging" {
		platformURL = DashboardStagingURL
	}
	return fmt.Sprintf("%s/settings/billing", platformURL)
}

// OpenBillingCmd represents the `meroxa open billing` command.
func OpenBillingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "billing",
		Short: "Open your billing page in a web browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			return browser.OpenURL(getBillingURL())
		},
	}
}
