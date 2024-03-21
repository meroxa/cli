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
	Image       string                 `json:"imageArchive"`
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
		builder.BuildCobraCommand(&Deploy{}),
		builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&Init{}),
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Open{}),
		builder.BuildCobraCommand(&Remove{}),
		builder.BuildCobraCommand(&Run{}),
	}
}

func RetrieveApplicationByNameOrID(ctx context.Context, client global.BasicClient, nameOrID string) (*Applications, error) {
	apps := Applications{}
	if nameOrID != "" {
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

	if apps.TotalItems == 0 {
		return nil, fmt.Errorf("no applications found")
	} else if apps.TotalItems > 1 {
		return nil, fmt.Errorf("multiple applications found")
	}

	return &apps, nil
}
