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

package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/environments"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type createResourceClient interface {
	CreateResource(ctx context.Context, resource *meroxa.CreateResourceInput) (*meroxa.Resource, error)
}

type Create struct {
	client createResourceClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		Type     string `long:"type"        short:""  usage:"resource type"        required:"true"`
		URL      string `long:"url"         short:"u" usage:"resource url"`
		Metadata string `long:"metadata"    short:"m" usage:"resource metadata"`

		// TODO: Add support to builder to create flags with an alias (--env | --environment)
		Environment string `long:"env" usage:"environment (name or UUID) where resource will be created"`

		// credentials
		Username      string `long:"username"    short:"" usage:"username"`
		Password      string `long:"password"    short:"" usage:"password"`
		CaCert        string `long:"ca-cert"     short:"" usage:"trusted certificates for verifying resource"`
		ClientCert    string `long:"client-cert" short:"" usage:"client certificate for authenticating to the resource"`
		ClientKey     string `long:"client-key"  short:"" usage:"client private key for authenticating to the resource"`
		SSL           bool   `long:"ssl"         short:"" usage:"use SSL"`
		SSHURL        string `long:"ssh-url"     short:"" usage:"SSH tunneling address"`
		SSHPrivateKey string `long:"ssh-private-key"     short:"" usage:"SSH tunneling private key"`
		Token         string `long:"token"     short:"" usage:"API Token"`
	}
}

var (
	_ builder.CommandWithDocs    = (*Create)(nil)
	_ builder.CommandWithArgs    = (*Create)(nil)
	_ builder.CommandWithFlags   = (*Create)(nil)
	_ builder.CommandWithClient  = (*Create)(nil)
	_ builder.CommandWithLogger  = (*Create)(nil)
	_ builder.CommandWithExecute = (*Create)(nil)
	_ builder.CommandWithAliases = (*Create)(nil)
)

func (c *Create) Usage() string {
	return "create [NAME] --type TYPE --url URL"
}

func (c *Create) Docs() builder.Docs {
	return builder.Docs{
		Short: "Add a resource to your Meroxa resource catalog",
		Long:  `Use the create command to add resources to your Meroxa resource catalog.`,

		// TODO: Provide example with `--env` once it's not behind a feature flag
		Example: `
$ meroxa resources create store --type postgres -u "$DATABASE_URL" --metadata '{"logical_replication":"true"}'
$ meroxa resources create datalake --type s3 -u "s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos"
$ meroxa resources create warehouse --type redshift -u "$REDSHIFT_URL"
$ meroxa resources create slack --type url -u "$WEBHOOK_URL"
$ meroxa resource create mysqldb \
    --type mysql \
    --url "mysql://$MYSQL_USER:$MYSQL_PASS@$MYSQL_URL:$MYSQL_PORT/$MYSQL_DB"

$ meroxa resource create mongo \
    --type mongodb \
    -u "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_URL:$MONGO_PORT"

$ meroxa resource create elasticsearch \
    --type elasticsearch \
    -u "https://$ES_USER:$ES_PASS@$ES_URL:$ES_PORT" \
    --metadata '{"index.prefix": "$ES_INDEX","incrementing.field.name": "$ES_INCREMENTING_FIELD"}'

$ meroxa resource create mybigquery \
    --type bigquery \
    -u "bigquery://$GCP_PROJECT_ID/$GCP_DATASET_NAME" \
    --client-key "$(cat $GCP_SERVICE_ACCOUNT_JSON_FILE)"

$ meroxa resource create snowflake \
    --type snowflakedb \
    -u "snowflake://$SNOWFLAKE_URL/meroxa_db/stream_data" \
    --username meroxa_user \
    --password $SNOWFLAKE_PRIVATE_KEY

$ meroxa resource create sourcedb \
	--type kafka \
	--url kafka+sasl+ssl://$KAFKA_USER:$KAFKA_PASS@<$BOOTSTRAP_SERVER>?sasl_mechanism=plain

$ meroxa resource create sourcedb \
	--type confluentcloud \
	--url kafka+sasl+ssl://$API_KEY:$API_SECRET@<$BOOTSTRAP_SERVER>?sasl_mechanism=plain`,
	}
}

func (c *Create) Client(client meroxa.Client) {
	c.client = client
}

func (c *Create) Logger(logger log.Logger) {
	c.logger = logger
}

func (c *Create) Aliases() []string {
	return []string{"add"}
}

func (c *Create) Flags() []builder.Flag {
	return builder.BuildFlags(&c.flags)
}

func (c *Create) ParseArgs(args []string) error {
	if len(args) > 0 {
		c.args.Name = args[0]
	}
	return nil
}

func (c *Create) Execute(ctx context.Context) error {
	var env string

	input := meroxa.CreateResourceInput{
		Type:     meroxa.ResourceTypeName(c.flags.Type),
		Name:     c.args.Name,
		URL:      c.flags.URL,
		Metadata: nil,
	}

	if c.flags.Type != string(meroxa.ResourceTypeNotion) && c.flags.URL == "" {
		return fmt.Errorf("required flag(s) \"url\" not set")
	}

	if c.flags.Environment != "" {
		err := builder.CheckCMDFeatureFlag(c, &environments.Environments{})
		if err != nil {
			return err
		}

		input.Environment = &meroxa.EntityIdentifier{}
		env = c.flags.Environment

		_, err = uuid.Parse(c.flags.Environment)

		if err == nil {
			input.Environment.UUID = c.flags.Environment
		} else {
			input.Environment.Name = c.flags.Environment
		}
	} else {
		env = string(meroxa.EnvironmentTypeCommon)
	}

	if c.hasCredentials() {
		input.Credentials = &meroxa.Credentials{
			Username:      c.flags.Username,
			Password:      c.flags.Password,
			CACert:        c.flags.CaCert,
			ClientCert:    c.flags.ClientCert,
			ClientCertKey: c.flags.ClientKey,
			UseSSL:        c.flags.SSL,
			Token:         c.flags.Token,
		}
	}

	if c.flags.Metadata != "" {
		err := json.Unmarshal([]byte(c.flags.Metadata), &input.Metadata)
		if err != nil {
			return fmt.Errorf("could not parse metadata: %w", err)
		}
	}

	if sshURL := c.flags.SSHURL; sshURL != "" {
		input.SSHTunnel = &meroxa.ResourceSSHTunnelInput{
			Address:    sshURL,
			PrivateKey: c.flags.SSHPrivateKey,
		}
	}

	c.logger.Infof(ctx, "Creating %q resource in %q environment...", input.Type, env)

	res, err := c.client.CreateResource(ctx, &input)
	if err != nil {
		return err
	}

	if tun := res.SSHTunnel; tun == nil {
		c.logger.Infof(ctx, "Resource %q is successfully created!", res.Name)
	} else {
		c.logger.Infof(ctx, "Resource %q is successfully created but is pending for validation!", res.Name)
		c.logger.Info(ctx, "Paste the following public key on your host:")
		c.logger.Info(ctx, tun.PublicKey)
		c.logger.Info(ctx, "Meroxa will try to connect to the resource for 60 minutes and send an email confirmation after a successful resource validation.") //nolint
	}

	c.logger.JSON(ctx, res)

	return nil
}

func (c *Create) hasCredentials() bool {
	return c.flags.Username != "" ||
		c.flags.Password != "" ||
		c.flags.CaCert != "" ||
		c.flags.ClientCert != "" ||
		c.flags.ClientKey != "" ||
		c.flags.Token != "" ||
		c.flags.SSL
}
