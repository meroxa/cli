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
	"os"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/prompt"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go"
)

var (
	_ builder.CommandWithDocs    = (*Remove)(nil)
	_ builder.CommandWithAliases = (*Remove)(nil)
	_ builder.CommandWithFlags   = (*Remove)(nil)
	_ builder.CommandWithArgs    = (*Remove)(nil)
	_ builder.CommandWithClient  = (*Remove)(nil)
	_ builder.CommandWithLogger  = (*Remove)(nil)
	_ builder.CommandWithExecute = (*Remove)(nil)
	_ builder.CommandWithPrompt  = (*Remove)(nil)
)

type removeEnvironmentClient interface {
	listEnvironmentsClient
	DeleteEnvironment(ctx context.Context, nameOrUUID string) (*meroxa.Environment, error)
}

type Remove struct {
	client removeEnvironmentClient
	logger log.Logger

	args struct {
		Name string
	}

	flags struct {
		NonInteractive bool `long:"no-input" usage:"skipping any prompts" default:"false"`
	}
}

func (r *Remove) Usage() string {
	return "remove NAME"
}

func (r *Remove) Docs() builder.Docs {
	return builder.Docs{
		Short: "Remove environment",
	}
}

func (r *Remove) Execute(ctx context.Context) error {
	r.logger.Infof(ctx, "Environment %q is being removed...", r.args.Name)

	e, err := r.client.DeleteEnvironment(ctx, r.args.Name)
	if err != nil {
		return err
	}

	r.logger.Infof(ctx, "Run `meroxa env describe %s` for status.", r.args.Name)
	r.logger.JSON(ctx, e)

	return nil
}

func (r *Remove) Prompt(ctx context.Context) error {
	if r.flags.NonInteractive {
		return nil
	}

	var list []string
	if r.args.Name == "" {
		envs, err := r.client.ListEnvironments(ctx)
		if err != nil {
			r.logger.Errorf(ctx, "cannot get envs: %+v\n", err)
			os.Exit(1)
		}
		for _, env := range envs {
			list = append(list, env.Name)
		}
		if len(list) == 0 {
			r.logger.Error(ctx, "no environment name provided, and no environments found")
			os.Exit(1)
		}
	}

	err := prompt.Show(ctx, []prompt.Prompt{
		prompt.SelectPrompt{
			Label:   "Select an environment to remove",
			Options: list,
			Value:   &r.args.Name,
			Skip:    r.args.Name != "",
		},
	})

	if err != nil {
		return err
	}
	if r.args.Name == "" {
		return builder.NewErrPromptAbort(
			errors.New("\nNo environment specified\n " +
				"please run \"meroxa help env remove\""),
		)
	}

	var remove bool
	_, err = prompt.BoolPrompt{
		Label: "Remove this environment",
		Value: &remove,
	}.Show(ctx)

	if err != nil {
		return err
	}

	if !remove {
		return builder.NewErrPromptAbort(
			errors.New("\nTo view all options for removing a Meroxa Environment,\n " +
				"please run \"meroxa help env remove\""),
		)
	}
	return nil
}

func (r *Remove) Logger(logger log.Logger) {
	r.logger = logger
}

func (r *Remove) Client(client *meroxa.Client) {
	r.client = client
}

func (r *Remove) Flags() []builder.Flag {
	return builder.BuildFlags(&r.flags)
}

func (r *Remove) ParseArgs(args []string) error {
	if len(args) >= 1 {
		r.args.Name = args[0]
	} else if len(args) < 1 && r.flags.NonInteractive == true {
		return fmt.Errorf("environment must be provided as an argument or interactively")
	}

	return nil
}

func (r *Remove) Aliases() []string {
	return []string{"rm", "delete"}
}
