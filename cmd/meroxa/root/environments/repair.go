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

package environments

import (
	"context"
	"errors"

	"github.com/meroxa/cli/utils/display"

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
	if len(args) < 1 {
		return errors.New("requires environment name or uuid")
	}
	r.args.NameOrUUID = args[0]
	return nil
}

func (r *Repair) Execute(ctx context.Context) error {
	environment, err := r.client.PerformActionOnEnvironment(ctx, r.args.NameOrUUID, &meroxa.RepairEnvironmentInput{Action: meroxa.EnvironmentActionRepair}) // nolint:lll
	if err != nil {
		return err
	}

	if environment.Status.State != meroxa.EnvironmentStatePreflightSuccess {
		details := display.EnvironmentPreflightTable(environment)
		r.logger.Errorf(ctx,
			"Environment %q could not be repaired because it failed the preflight checks\n%s\n",
			environment.Name,
			details)
	} else {
		r.logger.Infof(ctx,
			"Preflight checks have passed. Environment %q is being repaired. Run `meroxa env describe %s` for status",
			environment.Name,
			environment.Name)
	}

	r.logger.JSON(ctx, environment)
	return nil
}
