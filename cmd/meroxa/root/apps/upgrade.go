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
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine/javascript"
	turbinePY "github.com/meroxa/cli/cmd/meroxa/turbine/python"
	turbineRB "github.com/meroxa/cli/cmd/meroxa/turbine/ruby"
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
	return "upgrade [--path pwd]"
}

func (*Upgrade) Docs() builder.Docs {
	return builder.Docs{
		Short:   "Upgrade a Turbine Data Application",
		Example: `meroxa apps upgrade --path ~/code`,
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
		u.path, err = turbine.GetPath(u.flags.Path)
		u.logger.StartSpinner("\t", fmt.Sprintf("Fetching details of application in %q...", u.path))
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

	switch u.config.Language {
	case "go", turbine.GoLang:
		if u.turbineCLI == nil {
			u.turbineCLI = turbineGo.New(u.logger, u.path)
		}
	case "js", turbine.JavaScript, turbine.NodeJs:
		if u.turbineCLI == nil {
			u.turbineCLI = turbineJS.New(u.logger, u.path)
		}
	case "py", turbine.Python3, turbine.Python:
		if u.turbineCLI == nil {
			u.turbineCLI = turbinePY.New(u.logger, u.path)
		}
	case "rb", turbine.Ruby:
		if u.turbineCLI == nil {
			u.turbineCLI = turbineRB.New(u.logger, u.path)
		}
	default:
		return fmt.Errorf("language %q not supported. %s", u.config.Language, LanguageNotSupportedError)
	}
	vendor, _ := strconv.ParseBool(u.config.Vendor)
	if err = u.turbineCLI.Upgrade(vendor); err != nil {
		return err
	}

	u.logger.StartSpinner("\t", "Testing upgrades locally...")
	runOutput := ""
	buf := bytes.NewBufferString(runOutput)
	if u.run == nil {
		spinner := log.NewSpinnerLogger(buf)
		leveled := log.NewLeveledLogger(buf, log.Error)
		u.run = &Run{
			logger: log.New(leveled, nil, spinner),
			flags: struct {
				Path string `long:"path" usage:"path of application to run"`
			}{
				Path: u.path,
			},
		}
	}
	if err = u.run.Execute(ctx); err != nil {
		u.logger.Error(ctx, buf.String())
		u.logger.StopSpinnerWithStatus("Upgrades were not entirely successful."+
			" Fix any issues before adding and committing all upgrades.", log.Failed)
		return err
	}
	u.logger.StopSpinnerWithStatus("Tested upgrades locally successfully!", log.Successful)

	u.logger.Infof(ctx, "Your Turbine Application %s has been upgraded successfully!"+
		" To finish, add and commit the upgrades.", u.config.Name)
	return nil
}
