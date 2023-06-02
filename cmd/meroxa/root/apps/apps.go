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
	"fmt"

	turbineJava "github.com/meroxa/cli/cmd/meroxa/turbine/java"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine/javascript"
	turbinePY "github.com/meroxa/cli/cmd/meroxa/turbine/python"
	turbineRb "github.com/meroxa/cli/cmd/meroxa/turbine/ruby"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/spf13/cobra"
)

type Apps struct{}

var (
	_ builder.CommandWithDocs        = (*Apps)(nil)
	_ builder.CommandWithAliases     = (*Apps)(nil)
	_ builder.CommandWithSubCommands = (*Apps)(nil)
)

func (*Apps) Aliases() []string {
	return []string{"app"}
}

func (*Apps) Usage() string {
	return "apps"
}

func (*Apps) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage Turbine Data Applications",
	}
}

func (*Apps) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Deploy{}),
		builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&Init{}),
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Logs{}),
		builder.BuildCobraCommand(&Open{}),
		builder.BuildCobraCommand(&Remove{}),
		builder.BuildCobraCommand(&Run{}),
		builder.BuildCobraCommand(&Upgrade{}),
	}
}

// getTurbineCLIFromLanguage will return the appropriate turbine.CLI based on language.
func getTurbineCLIFromLanguage(logger log.Logger, lang ir.Lang, path string) (turbine.CLI, error) {
	switch lang {
	case "go", turbine.GoLang:
		return turbineGo.New(logger, path), nil
	case "js", turbine.JavaScript, turbine.NodeJs:
		return turbineJS.New(logger, path), nil
	case "py", turbine.Python3, turbine.Python:
		return turbinePY.New(logger, path), nil
	case "rb", turbine.Ruby:
		return turbineRb.New(logger, path), nil
	case turbine.Java:
		return turbineJava.New(logger, path), nil
	}
	return nil, newLangUnsupportedError(lang)
}

type addHeader interface {
	AddHeader(key, value string)
}

func addTurbineHeaders(c addHeader, lang ir.Lang, version string) {
	c.AddHeader("Meroxa-CLI-App-Lang", string(lang))
	if lang == ir.JavaScript {
		version = fmt.Sprintf("%s:cli%s", version, turbineJS.TurbineJSVersion)
	}
	c.AddHeader("Meroxa-CLI-App-Version", version)
}
