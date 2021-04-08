package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

// ApiCmd represents the `meroxa api` command
func ApiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api METHOD PATH [body]",
		Short: "Invoke Meroxa API",
		Args:  cobra.MinimumNArgs(2),
		Example: `
meroxa api GET /v1/endpoints
meroxa api POST /v1/endpoints '{"protocol": "HTTP", "stream": "resource-2-499379-public.accounts", "name": "1234"}'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client()
			if err != nil {
				return err
			}

			var (
				method = args[0]
				path   = args[1]
				body   string
			)

			if len(args) > 2 {
				body = args[2]
			}

			resp, err := c.MakeRequestString(context.Background(), method, path, body)
			if err != nil {
				return err
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, b, "", "\t"); err != nil {
				prettyJSON.Write(b)
			}

			fmt.Printf("> %s %s\n", method, path)
			fmt.Printf("< %s %s\n", resp.Status, resp.Proto)
			for k, v := range resp.Header {
				fmt.Printf("< %s %s\n", k, strings.Join(v, " "))
			}
			fmt.Printf(prettyJSON.String())

			return nil
		},
	}
}
