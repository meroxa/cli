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
func OpenBillingCmd() *cobra.Command {
	return &cobra.Command{
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
}