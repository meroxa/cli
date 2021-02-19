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

package cmd

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"os"
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
	return fmt.Sprintf("%s/account/billing", platformURL)
}

// openBillingCmd represents the billing command
var openBillingCmd = &cobra.Command{
	Use:   "billing",
	Short: "Open your billing page in a web browser",
	RunE: func(cmd *cobra.Command, args []string) error {

		err := browser.OpenURL(getBillingURL())

		if err != nil {
			return err
		}

		return nil
	},
}

// openCmd represents the billing command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open in a web browser",
}

func init() {
	openCmd.AddCommand(openBillingCmd)
	RootCmd.AddCommand(openCmd)
}
