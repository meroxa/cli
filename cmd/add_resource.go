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
	"context"
	"encoding/json"
	"fmt"

	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type AddResourceClient interface {
	CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error)
}

type AddResource struct {
	Name       string
	Type       string `mapstructure:"type"`
	Url        string `mapstructure:"url"`
	Metadata   string `mapstructure:"metadata"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	CaCert     string `mapstructure:"ca-cert"`
	ClientCert string `mapstructure:"client-cert"`
	ClientKey  string `mapstructure:"client-key"`
	Ssl        bool   `mapstructure:"ssl"`
	cfg        *viper.Viper
}

func (ar *AddResource) setArgs(args []string) error {
	if len(args) > 0 {
		ar.Name = args[0]
	}

	return nil
}

func (ar *AddResource) setFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("type", "", "", "resource type")
	cmd.MarkFlagRequired("type")

	cmd.Flags().StringP("url", "u", "", "resource url")
	cmd.MarkFlagRequired("url")

	cmd.Flags().StringP("username", "", "", "username")
	cmd.Flags().StringP("password", "", "", "passsword")
	cmd.Flags().StringP("ca-cert", "", "", "trusted certificates for verifying resource")
	cmd.Flags().StringP("client-cert", "", "", "client certificate for authenticating to the resource")
	cmd.Flags().StringP("client-key", "", "", "client private key for authenticating to the resource")
	cmd.Flags().BoolP("ssl", "", false, "use SSL")

	cmd.Flags().StringP("metadata", "m", "", "resource metadata")

	viperBindFlags(cmd, ar.cfg)
}

func viperBindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flagName := flag.Name
		if flagName != "config" && flagName != "help" {
			if err := v.BindPFlag(flagName, flag); err != nil {
				panic(fmt.Errorf("error binding flag '%s': %w", flagName, err).Error())
			}
		}
	})
}

func (ar *AddResource) execute(ctx context.Context, c AddResourceClient, res meroxa.CreateResourceInput) (*meroxa.Resource, error) {
	if !flagRootOutputJSON {
		fmt.Printf("Adding %s resource...\n", res.Type)
	}

	var err error

	res.Credentials = &meroxa.Credentials{
		Username:      ar.Username,
		Password:      ar.Password,
		CACert:        ar.CaCert,
		ClientCert:    ar.ClientCert,
		ClientCertKey: ar.ClientKey,
		UseSSL:        ar.Ssl,
	}

	if ar.Metadata != "" {
		var metadata map[string]interface{}
		err = json.Unmarshal([]byte(ar.Metadata), &metadata)
		if err != nil {
			return nil, err
		}

		res.Metadata = metadata
	}

	resource, err := c.CreateResource(ctx, &res)
	return resource, err
}

func (ar *AddResource) output(res *meroxa.Resource) {
	if flagRootOutputJSON {
		utils.JSONPrint(res)
	} else {
		fmt.Printf("%s resource with name %s successfully added!\n", res.Type, res.Name)
	}
}

// AddResourceCmd represents the `meroxa add resource` command
func (ar *AddResource) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource [NAME] --type TYPE",
		Short: "Add a resource to your Meroxa resource catalog",
		Long:  `Use the add command to add resources to your Meroxa resource catalog.`,
		Example: "\n" +
			"meroxa add resource store --type postgres -u $DATABASE_URL --metadata '{\"logical_replication\":true}'\n" +
			"meroxa add resource datalake --type s3 -u \"s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos\"\n" +
			"meroxa add resource warehouse --type redshift -u $REDSHIFT_URL\n" +
			"meroxa add resource slack --type url -u $WEBHOOK_URL\n",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ar.setArgs(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ar.cfg.Unmarshal(ar); err != nil {
				return err
			}

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, clientTimeOut)
			defer cancel()

			c, err := client()

			if err != nil {
				return err
			}

			ri := meroxa.CreateResourceInput{
				Type:     ar.Type,
				Name:     ar.Name,
				URL:      ar.Url,
				Metadata: nil,
			}

			res, err := ar.execute(ctx, c, ri)

			if err != nil {
				return err
			}

			ar.output(res)

			return nil
		},
	}

	ar.setFlags(cmd)

	return cmd
}
