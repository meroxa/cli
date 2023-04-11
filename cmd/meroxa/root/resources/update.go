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
	"errors"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

type updateResourceClient interface {
	GetResourceByNameOrID(ctx context.Context, nameOrID string) (*meroxa.Resource, error)
	UpdateResource(ctx context.Context, nameOrID string, resourceToUpdate *meroxa.UpdateResourceInput) (*meroxa.Resource, error)
}

type Update struct {
	client updateResourceClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		URL      string `long:"url"         short:"u" usage:"new resource url"`
		Metadata string `long:"metadata"    short:"m" usage:"new resource metadata"`
		Name     string `long:"name"        usage:"new resource name"`

		// credentials
		Username   string `long:"username"    short:"" usage:"username"`
		Password   string `long:"password"    short:"" usage:"password"`
		CaCert     string `long:"ca-cert"     short:"" usage:"trusted certificates for verifying resource"`
		ClientCert string `long:"client-cert" short:"" usage:"client certificate for authenticating to the resource"`
		ClientKey  string `long:"client-key"  short:"" usage:"client private key for authenticating to the resource"`
		SSL        bool   `long:"ssl"         short:"" usage:"use SSL"`
		SSHURL     string `long:"ssh-url"     short:"" usage:"SSH tunneling address"`
		Token      string `long:"token"       short:"" usage:"API Token"`
	}
}

func (u *Update) Usage() string {
	return "update NAME"
}

func (u *Update) Docs() builder.Docs {
	return builder.Docs{
		Short: "Update a resource",
		Long:  "Use the update command to update various Meroxa resources.",
	}
}

func (u *Update) Execute(ctx context.Context) error {
	// TODO: Implement something like dependent flags in Builder
	if u.flags.Name == "" && u.flags.URL == "" && u.flags.Metadata == "" && !u.isUpdatingCredentials() {
		return errors.New("requires either `--name`, `--url`, `--metadata` or one of the credential flags")
	}

	r, err := u.client.GetResourceByNameOrID(ctx, u.args.Name)
	if err != nil {
		return err
	}

	u.logger.Infof(ctx, "Updating resource %q...", u.args.Name)

	res := &meroxa.UpdateResourceInput{}

	// If name was provided, update it
	if u.flags.Name != "" {
		res.Name = u.flags.Name
	}

	// If url was provided, update it
	if u.flags.URL != "" {
		u.processURLFlag(ctx, string(r.Type))
		res.URL = u.flags.URL
	}

	// If metadata was provided, update it
	if u.flags.Metadata != "" {
		var metadata map[string]interface{}
		err := json.Unmarshal([]byte(u.flags.Metadata), &metadata)
		if err != nil {
			return fmt.Errorf("can't parse metadata: %w", err)
		}
		res.Metadata = metadata
	}

	if sshURL := u.flags.SSHURL; sshURL != "" {
		res.SSHTunnel = &meroxa.ResourceSSHTunnelInput{
			Address: sshURL,
		}
	}

	// If any of the credential values are being updated
	if u.isUpdatingCredentials() {
		res.Credentials = &meroxa.Credentials{
			Username:      u.flags.Username,
			Password:      u.flags.Password,
			CACert:        u.flags.CaCert,
			ClientCert:    u.flags.ClientCert,
			ClientCertKey: u.flags.ClientKey,
			UseSSL:        u.flags.SSL,
			Token:         u.flags.Token,
		}
	}

	r, err = u.client.UpdateResource(ctx, u.args.Name, res)

	if err != nil {
		return err
	}

	if tun := r.SSHTunnel; tun == nil {
		u.logger.Infof(ctx, "Resource %q is successfully updated!", r.Name)
	} else {
		u.logger.Infof(ctx, "Resource %q is successfully updated but is pending for validation!", r.Name)
		u.logger.Info(ctx, "Meroxa will try to connect to the resource for 60 minutes and send an email confirmation after a successful resource validation.") //nolint
	}

	u.logger.JSON(ctx, r)
	return nil
}

func (u *Update) Flags() []builder.Flag {
	return builder.BuildFlags(&u.flags)
}

func (u *Update) Logger(logger log.Logger) {
	u.logger = logger
}

func (u *Update) Client(client meroxa.Client) {
	u.client = client
}

func (u *Update) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires resource name")
	}

	u.args.Name = args[0]
	return nil
}

var (
	_ builder.CommandWithDocs    = (*Update)(nil)
	_ builder.CommandWithArgs    = (*Update)(nil)
	_ builder.CommandWithFlags   = (*Update)(nil)
	_ builder.CommandWithClient  = (*Update)(nil)
	_ builder.CommandWithLogger  = (*Update)(nil)
	_ builder.CommandWithExecute = (*Update)(nil)
)

func (u *Update) isUpdatingCredentials() bool {
	return u.flags.Username != "" ||
		u.flags.Password != "" ||
		u.flags.CaCert != "" ||
		u.flags.ClientCert != "" ||
		u.flags.ClientKey != "" ||
		u.flags.Token != "" ||
		u.flags.SSL
}

func (u *Update) processURLFlag(ctx context.Context, rt string) {
	if rt == string(meroxa.ResourceTypeNotion) {
		url := u.flags.URL
		u.flags.URL = ""
		if url != "" && url != defaultNotionUrl {
			u.logger.Warnf(ctx, "Ignoring API URL override (%s) for Notion resource configuration.", url)
		}
	} else if rt == string(meroxa.ResourceTypeSpireMaritimeAIS) {
		url := u.flags.URL
		u.flags.URL = ""
		if url != "" && url != defaultSpireMaritimeAisUrl {
			u.logger.Warnf(ctx, "Ignoring API URL override (%s) for Spire Maritime AIS resource configuration.", url)
		}
	}
}
