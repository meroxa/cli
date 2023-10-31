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
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	turbineGo "github.com/meroxa/cli/cmd/meroxa/turbine/golang"
	turbineJS "github.com/meroxa/cli/cmd/meroxa/turbine/javascript"
	turbinePY "github.com/meroxa/cli/cmd/meroxa/turbine/python"
	turbineRb "github.com/meroxa/cli/cmd/meroxa/turbine/ruby"
	pb "github.com/pocketbase/pocketbase/tools/types"

	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils/display"
	"github.com/meroxa/turbine-core/pkg/ir"
	"github.com/spf13/cobra"
)

type ApplicationState string

const (
	ApplicationStateInitialized ApplicationState = "initialized"
	ApplicationStateDeploying   ApplicationState = "deploying"
	ApplicationStatePending     ApplicationState = "pending"
	ApplicationStateRunning     ApplicationState = "running"
	ApplicationStateDegraded    ApplicationState = "degraded"
	ApplicationStateFailed      ApplicationState = "failed"

	collectionName = "apps"
)

var displayDetails = display.Details{
	"Name":        "name",
	"State":       "state",
	"SpecVersion": "specVersion",
	"Created":     "created",
	"Updated":     "updated",
}

// Application represents the Meroxa Application type within the Meroxa API.
type Application struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	State       ApplicationState       `json:"state"`
	Spec        map[string]interface{} `json:"spec"`
	SpecVersion string                 `json:"specVersion"`
	Created     AppTime                `json:"created"`
	Updated     AppTime                `json:"updated"`
}

type Applications struct {
	Page       int           `json:"page"`
	PerPage    int           `json:"perPage"`
	TotalItems int           `json:"totalItems"`
	TotalPages int           `json:"totalPages"`
	Items      []Application `json:"items"`
}

type AppTime struct {
	time.Time
}

func (at *AppTime) UnmarshalJSON(b []byte) error {
	appTime, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	dt, err := pb.ParseDateTime(appTime) // time.Parse(pb.DefaultDateLayout, appTime)
	if err != nil {
		fmt.Println(err)
		return err
	}
	at.Time = dt.Time()
	return nil
}

func (at *AppTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(at.Time)
}

func (at *AppTime) Format(s string) string {
	t := at.Time
	return t.Format(s)
}

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
		builder.BuildCobraCommand(&Open{}),
		builder.BuildCobraCommand(&Remove{}),
		builder.BuildCobraCommand(&Run{}),
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

func (a Applications) RetrieveApplicationID(ctx context.Context, client global.BasicClient, nameOrID, path string) (*Applications, error) {
	var getPath string
	apps := Applications{}
	if path != "" {
		var err error
		if getPath, err = turbine.GetPath(path); err != nil {
			return nil, err
		}

		config, err := turbine.ReadConfigFile(getPath)
		if err != nil {
			return nil, err
		}

		a := &url.Values{}
		a.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", config.Name, config.Name))

		response, err := client.CollectionRequest(ctx, "GET", collectionName, "", nil, *a)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(response.Body).Decode(&apps)
		if err != nil {
			return nil, err
		}
	} else if nameOrID != "" {
		a := &url.Values{}
		a.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", nameOrID, nameOrID))

		response, err := client.CollectionRequest(ctx, "GET", collectionName, "", nil, *a)
		if err != nil {
			return nil, err
		}

		err = json.NewDecoder(response.Body).Decode(&apps)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("supply either ID/Name argument or --path flag")
	}
	return &apps, nil
}
