package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	pb "github.com/pocketbase/pocketbase/tools/types"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/utils/display"
	"github.com/spf13/cobra"
)

const (
	collectionName = "secrets"
)

type ListSecrets struct {
	Page       int       `json:"page"`
	PerPage    int       `json:"perPage"`
	TotalItems int       `json:"totalItems"`
	TotalPages int       `json:"totalPages"`
	Items      []Secrets `json:"items"`
}

type Secrets struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Data           map[string]interface{} `json:"data"`
	Created        SecretTime             `json:"created"`
	Updated        SecretTime             `json:"updated"`
	CollectionID   string                 `json:"collectionId"`
	CollectionName string                 `json:"collectionName"`
}

var displayDetails = display.Details{
	"ID":      "id",
	"Name":    "name",
	"Data":    "data",
	"Created": "created",
	"Updated": "updated",
}

var (
	_ builder.CommandWithDocs        = (*Secrets)(nil)
	_ builder.CommandWithAliases     = (*Secrets)(nil)
	_ builder.CommandWithSubCommands = (*Secrets)(nil)
)

type SecretTime struct {
	time.Time
}

func (at *SecretTime) UnmarshalJSON(b []byte) error {
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

func (at *SecretTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(at.Time)
}

func (at *SecretTime) Format(s string) string {
	t := at.Time
	return t.Format(s)
}

func (*Secrets) Aliases() []string {
	return []string{"secrets"}
}

func (*Secrets) Usage() string {
	return "secrets"
}

func (*Secrets) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage Turbine Data Applications",
	}
}

func (*Secrets) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Create{}),
		builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&List{}),
		builder.BuildCobraCommand(&Remove{}),
	}
}

func RetrieveSecretsID(ctx context.Context, client global.BasicClient, nameOrID string) (*ListSecrets, error) {
	getSecrets := &ListSecrets{}

	a := &url.Values{}
	a.Add("filter", fmt.Sprintf("(id='%s' || name='%s')", nameOrID, nameOrID))

	response, err := client.CollectionRequest(ctx, "GET", collectionName, "", nil, *a)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&getSecrets)
	if err != nil {
		return nil, err
	}

	if getSecrets.TotalItems == 0 {
		return nil, fmt.Errorf("secret %q not found", nameOrID)
	} else if getSecrets.TotalItems > 1 {
		return nil, fmt.Errorf("multiple secrets found with name %q", nameOrID)
	}

	return getSecrets, nil
}
