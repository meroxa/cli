/*
Copyright © 2021 Meroxa Inc
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

package environments

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Repair)(nil)
	_ builder.CommandWithArgs    = (*Repair)(nil)
	_ builder.CommandWithClient  = (*Repair)(nil)
	_ builder.CommandWithLogger  = (*Repair)(nil)
	_ builder.CommandWithExecute = (*Repair)(nil)
	_ builder.CommandWithHidden  = (*Repair)(nil)
)

type repairEnvironmentClient interface {
	PerformActionOnEnvironment(ctx context.Context, nameOrUUID string, body *meroxa.RepairEnvironmentInput) (*meroxa.Environment, error)
}

type Repair struct {
	client repairEnvironmentClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}
}

func (r *Repair) Hidden() bool {
	return true
}

func (r *Repair) Usage() string {
	return "repair NAMEorUUID"
}

func (r *Repair) Docs() builder.Docs {
	return builder.Docs{
		Short: "Repair environment",
		Long:  `Repair any environment that is in one of the following states: provisioning_error, deprovisioning_error, repairing_error.`,
	}
}

func (r *Repair) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Repair) Client(client meroxa.Client) {
	r.client = client
}

func (r *Repair) ParseArgs(args []string) error {
	if len(args) > 0 {
		r.args.NameOrUUID = args[0]
	}
	return nil
}

func (r *Repair) Execute(ctx context.Context) error {
	rr, err := r.client.PerformActionOnEnvironment(ctx, r.args.NameOrUUID, &meroxa.RepairEnvironmentInput{Action: "repair"})
	if err != nil {
		return err
	}

	r.logger.Infof(ctx, `The repairment of your environment %q is now in progress and your environment will be up and running soon.`, r.args.NameOrUUID) // nolint:lll
	r.logger.Infof(ctx, `Run "meroxa env describe %s" for status.`, r.args.NameOrUUID)
	r.logger.JSON(ctx, rr)

	return nil
}
