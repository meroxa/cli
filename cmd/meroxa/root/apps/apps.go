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
	pb "github.com/pocketbase/pocketbase/tools/types"

	"github.com/meroxa/cli/utils/display"
	"github.com/spf13/cobra"
)

type ApplicationState string

const (
	AppProvisioningState   ApplicationState = "provisioning"
	AppProvisionedState    ApplicationState = "provisioned"
	AppDeprovisioningState ApplicationState = "deprovisioning"
	AppDeprovisionedState  ApplicationState = "deprovisioned"
	AppDegradedState       ApplicationState = "degraded"

	deploymentCollection  = "conduit_deployments"
	applicationCollection = "conduit_apps"
)

var displayDetails = display.Details{
	"Name":              "name",
	"State":             "state",
	"ApplicationSpec":   "stream_tech",
	"Config":            "config",
	"PipelineFilenames": "pipelines_filenames",
	// "PipelineEnriched":  "pipeline_enriched",
	// "PipelineOriginal":  "pipeline_original",
	"Created": "created",
	"Updated": "updated",
}

type Deployment struct {
	ID              string `json:"id,omitempty"`
	Archive         string `json:"archive"`
	State           string `json:"state,omitempty"`
	ApplicationSpec string `json:"app_spec,omitempty"`

	Created            AppTime `json:"created,omitempty"`
	Updated            AppTime `json:"updated,omitempty"`
	ProcessorPlugins   string  `json:"processors_plugins,omitempty"`
	ProcessorFilenames string  `json:"processors_filenames,omitempty"`
	PipelineFilenames  string  `json:"pipeline_filenames,omitempty"`
}

type Deployments struct {
	Page       int          `json:"page"`
	PerPage    int          `json:"perPage"`
	TotalItems int          `json:"totalItems"`
	TotalPages int          `json:"totalPages"`
	Items      []Deployment `json:"items"`
}

// Application represents the Meroxa Application type within the Meroxa API.
type Application struct {
	ID                string   `json:"id"`
	DeploymentID      []string `json:"deployment_id"`
	Name              string   `json:"name"`
	State             string   `json:"state"`
	ApplicationSpec   string   `json:"stream_tech"`
	Config            string   `json:"config"`
	PipelineFilenames string   `json:"pipeline_filenames"`
	PipelineEnriched  string   `json:"pipeline_enriched"`
	PipelineOriginal  string   `json:"pipeline_original"`

	Created AppTime `json:"created"`
	Updated AppTime `json:"updated"`
	Archive string  `json:"archive"`
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
		Short: "Manage Conduit Data Applications",
	}
}

func (*Apps) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		// TODO - commenting out run and init until implemented
		builder.BuildCobraCommand(&Deploy{}),
		builder.BuildCobraCommand(&Describe{}),
		// builder.BuildCobraCommand(&Init{}),
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Open{}),
		builder.BuildCobraCommand(&Remove{}),
		// builder.BuildCobraCommand(&Run{}),
	}
}

func RetrieveApplicationByNameOrID(ctx context.Context, client global.BasicClient, nameOrID, path string) (*Applications, error) {
	apps := Applications{}
	if nameOrID != "" {
		a := &url.Values{}
		a.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", nameOrID, nameOrID))

		response, err := client.CollectionRequest(ctx, "GET", applicationCollection, "", nil, *a)
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

	if apps.TotalItems == 0 {
		return nil, fmt.Errorf("no applications found")
	} else if apps.TotalItems > 1 {
		return nil, fmt.Errorf("multiple applications found")
	}

	return &apps, nil
}
