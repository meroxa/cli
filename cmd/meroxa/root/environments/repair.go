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
	"errors"
	"fmt"

	"github.com/meroxa/cli/utils"

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
	if len(args) < 1 {
		return errors.New("requires environment name or uuid")
	}
	r.args.NameOrUUID = args[0]
	return nil
}

func (r *Repair) Execute(ctx context.Context) error {
	rr, err := r.client.PerformActionOnEnvironment(ctx, r.args.NameOrUUID, &meroxa.RepairEnvironmentInput{Action: meroxa.EnvironmentActionRepair}) // nolint:lll
	if err != nil {
		return err
	}

	state := rr.Status.State
	name := rr.Name
	if rr != nil && state == meroxa.EnvironmentStatePreflightError {
		log := fmt.Sprintf("Environment %q could not be repaired because it failed the preflight checks.", name)
		details := utils.PrettyString(rr.Status.PreflightDetails)
		if details != "" && details != "null" {
			log += fmt.Sprintf("\n%s\n", details)
		}
		r.logger.Errorf(ctx, log)
	} else if rr != nil && (state != meroxa.EnvironmentStateRepairing && state != meroxa.EnvironmentStateReady) {
		log := fmt.Sprintf("Environment %q could not be repaired.", r.args.NameOrUUID)
		details := utils.PrettyString(rr.Status.Details)
		if details != "" && details != "null" {
			log += fmt.Sprintf("\n%s\n", details)
		}
		switch state {
		case meroxa.EnvironmentStateRepairingError:
		case meroxa.EnvironmentStateUpdatingError:
		case meroxa.EnvironmentStateProvisioningError:
		case meroxa.EnvironmentStateDeprovisioningError:
			r.logger.Infof(ctx, log)
		}
	} else {
		r.logger.Infof(ctx, `The repairment of your environment %q is now in progress and your environment will be up and running soon.`, r.args.NameOrUUID) // nolint:lll
	}
	r.logger.Infof(ctx, `Run "meroxa env describe %s" for status.`, r.args.NameOrUUID)
	r.logger.JSON(ctx, rr)

	return nil
}
