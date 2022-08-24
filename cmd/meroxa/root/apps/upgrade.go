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

package apps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine/javascript"
	turbinePY "github.com/meroxa/cli/cmd/meroxa/turbine/python"
	"github.com/meroxa/cli/log"
)

type Upgrade struct {
	logger     log.Logger
	turbineCLI turbine.CLI
	run        builder.CommandWithExecute
	path       string
	config     *turbine.AppConfig

	flags struct {
		Path string `long:"path" usage:"path where application exists (current directory as default)"`
	}
}

var (
	_ builder.CommandWithDocs    = (*Upgrade)(nil)
	_ builder.CommandWithFlags   = (*Upgrade)(nil)
	_ builder.CommandWithExecute = (*Upgrade)(nil)
	_ builder.CommandWithLogger  = (*Upgrade)(nil)
)

func (*Upgrade) Usage() string {
	return "upgrade [APP_NAME] [--path pwd]"
}

func (*Upgrade) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Upgrade a Turbine Data Application",
		Example: `meroxa apps upgrade my-app --path ~/code`,
		Beta:    true,
	}
}

func (u *Upgrade) Logger(logger log.Logger) {
	u.logger = logger
}

func (u *Upgrade) Flags() []builder.Flag {
	return builder.BuildFlags(&u.flags)
}

func (u *Upgrade) Execute(ctx context.Context) error {
	var err error
	if u.config == nil {
		u.logger.StartSpinner("\t", fmt.Sprintf(" Fetching details of application in %q...", u.path))
		u.path, err = turbine.GetPath(u.flags.Path)
		if err != nil {
			u.logger.StopSpinnerWithStatus("\t", log.Failed)
			return err
		}

		var config turbine.AppConfig
		config, err = turbine.ReadConfigFile(u.path)
		u.config = &config
		if err != nil {
			u.logger.StopSpinnerWithStatus("\t", log.Failed)
			return err
		}
		u.logger.StopSpinnerWithStatus(fmt.Sprintf("Determined the details of the %q Application", u.config.Name), log.Successful)
	}

	lang := u.config.Language
	vendor, _ := strconv.ParseBool(u.config.Vendor)
	switch lang {
	case "go", GoLang:
		if u.turbineCLI == nil {
			u.turbineCLI = turbineGo.New(u.logger, u.path)
		}
		err = u.turbineCLI.Upgrade(vendor)
	case "js", JavaScript, NodeJs:
		if u.turbineCLI == nil {
			u.turbineCLI = turbineJS.New(u.logger, u.path)
		}
		err = u.turbineCLI.Upgrade(vendor)
	case "py", Python3, Python:
		if u.turbineCLI == nil {
			u.turbineCLI = turbinePY.New(u.logger, u.path)
		}
		err = u.turbineCLI.Upgrade(vendor)
	default:
		return fmt.Errorf("language %q not supported. %s", lang, LanguageNotSupportedError)
	}
	if err != nil {
		return err
	}

	u.logger.StartSpinner("\t", " Testing upgrades locally...")
	if u.run == nil {
		u.run = &Run{
			logger: log.NewWithDevNull(),
			flags: struct {
				Path string `long:"path" usage:"path of application to run"`
			}{
				Path: u.path,
			},
		}
	}
	err = u.run.Execute(ctx)
	if err != nil {
		u.logger.StopSpinnerWithStatus("Upgrades were not entirely successful."+
			" Fix any issues before adding and committing all upgrades.", log.Failed)
		return err
	}
	u.logger.StopSpinnerWithStatus("Tested upgrades locally successfully!", log.Successful)

	u.logger.Infof(ctx, "Your Turbine Application %s has been upgraded successfully!"+
		" To finish, add and commit the upgrades.", u.config.Name)
	return nil
}
