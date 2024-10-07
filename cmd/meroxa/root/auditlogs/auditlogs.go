package auditlogs

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
	collectionName = "auditlog"
)

type ListAuditlog struct {
	Page       int        `json:"page"`
	PerPage    int        `json:"perPage"`
	TotalItems int        `json:"totalItems"`
	TotalPages int        `json:"totalPages"`
	Items      []Auditlog `json:"items"`
}

type Auditlog struct {
	ID         string       `json:"id"`
	Collection string       `json:"collection"`
	Event      string       `json:"event"`
	UserID     string       `json:"user_id"`
	User       string       `json:"user"`
	Created    auditlogTime `json:"created"`
	Updated    auditlogTime `json:"updated"`
}

var displayDetails = display.Details{
	"ID":         "id",
	"Event":      "event",
	"Collection": "collection",
	"UserID":     "user_id",
	"User":       "user",
	"Created":    "created",
	"Updated":    "updated",
}

var (
	_ builder.CommandWithDocs        = (*Auditlog)(nil)
	_ builder.CommandWithAliases     = (*Auditlog)(nil)
	_ builder.CommandWithSubCommands = (*Auditlog)(nil)
)

type auditlogTime struct {
	time.Time
}

func (at *auditlogTime) UnmarshalJSON(b []byte) error {
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

func (at *auditlogTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(at.Time)
}

func (at *auditlogTime) Format(s string) string {
	t := at.Time
	return t.Format(s)
}

func (*Auditlog) Aliases() []string {
	return []string{"auditlog"}
}

func (*Auditlog) Usage() string {
	return "auditlog"
}

func (*Auditlog) Docs() builder.Docs {
	return builder.Docs{
		Short: "Manage Conduit Application Auditlog",
	}
}

func (*Auditlog) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		//	builder.BuildCobraCommand(&Create{}),
		//	builder.BuildCobraCommand(&Describe{}),
		builder.BuildCobraCommand(&List{}),
		//	builder.BuildCobraCommand(&Remove{}),
	}
}

func RetrieveAuditlogID(ctx context.Context, client global.BasicClient, nameOrID string) (*ListAuditlog, error) {
	getAuditlog := &ListAuditlog{}

	a := &url.Values{}

	response, err := client.CollectionRequest(ctx, "GET", collectionName, "", nil, *a)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&getAuditlog)
	if err != nil {
		return nil, err
	}

	if getAuditlog.TotalItems == 0 {
		return nil, fmt.Errorf("auditlog %q not found", nameOrID)
	} else if getAuditlog.TotalItems > 1 {
		return nil, fmt.Errorf("multiple Auditlog found with name %q", nameOrID)
	}

	return getAuditlog, nil
}
