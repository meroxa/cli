package fixtures

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/turbine"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
)

var (
	_ builder.CommandWithDocs    = (*Fetch)(nil)
	_ builder.CommandWithArgs    = (*Fetch)(nil)
	_ builder.CommandWithFlags   = (*Fetch)(nil)
	_ builder.CommandWithClient  = (*Fetch)(nil)
	_ builder.CommandWithLogger  = (*Fetch)(nil)
	_ builder.CommandWithExecute = (*Fetch)(nil)
)

type introspectResourceClient interface {
	IntrospectResource(ctx context.Context, nameOrUUID string) (*meroxa.ResourceIntrospection, error)
}

type Fetch struct {
	client introspectResourceClient
	logger log.Logger

	args struct {
		NameOrUUID string
	}

	flags struct {
		Collections string `long:"collections" short:"c" usage:"fetch --collections table1,table2"`
		Path        string `long:"path" short:"p" usage:"fetch --path ./foo/bar"`
	}
}

func (f *Fetch) Usage() string {
	return "fetch [NAMEorUUID]"
}

func (f *Fetch) Docs() builder.Docs {
	return builder.Docs{
		Short: "Fetch fixtures",
		Long: "Fetch fixtures retrieves sample data records from the provided resource and makes them " +
			"available in the \"/fixtures\" sub-directory",
	}
}

func (f *Fetch) Flags() []builder.Flag {
	return builder.BuildFlags(&f.flags)
}

func (f *Fetch) Execute(ctx context.Context) error {
	resourceName := f.args.NameOrUUID
	f.logger.Infof(ctx, "Fetching fixtures for %s...", resourceName)
	ri, err := f.client.IntrospectResource(ctx, resourceName)
	if err != nil {
		return err
	}

	var samples Samples
	if f.flags.Collections != "" {
		cols := strings.Split(f.flags.Collections, ",")
		samples = filterCollections(cols, ri.Samples)
	} else {
		samples = ri.Samples
	}

	var pwd string
	if f.flags.Path == "" {
		pwd, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		pwd = f.flags.Path
	}

	appConfig, err := turbine.ReadConfigFile(pwd)
	if err != nil {
		return err
	}

	// iterate through resources and write samples to the paths listed
	for r, c := range appConfig.Resources {
		if r == resourceName {
			pathToFixtureFile := filepath.Join(pwd, c)
			f.logger.Infof(ctx, "Writing fixtures to %s", pathToFixtureFile)
			err := writeFixtures(pathToFixtureFile, samples)
			if err != nil {
				return err
			}
		}
	}

	f.logger.Info(ctx, "Successfully fetched fixtures")

	return nil
}

func (f *Fetch) Client(client meroxa.Client) {
	f.client = client
}

func (f *Fetch) Logger(logger log.Logger) {
	f.logger = logger
}

func (f *Fetch) ParseArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("requires function name")
	}

	f.args.NameOrUUID = args[0]
	return nil
}

type Samples map[string][]string

func filterCollections(collections []string, samples Samples) Samples {
	filteredCollections := make(map[string][]string)
	for _, c := range collections {

		if sample, ok := samples[c]; ok {
			filteredCollections[c] = sample
		}
	}
	return filteredCollections
}

func writeFixtures(path string, samples Samples) error {
	// this is annoying but needed to unquote the nested JSON within each collection as each record
	// in the collection array is JSON formatted already.
	smap := make(map[string]interface{})
	for col, recs := range samples {
		smap[col] = []map[string]interface{}{}
		for _, r := range recs {
			var m map[string]interface{}
			err := json.Unmarshal([]byte(r), &m)
			if err != nil {
				return err
			}
			smap[col] = append(smap[col].([]map[string]interface{}), m)
		}
	}

	jSamples, err := json.MarshalIndent(smap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, jSamples, 0644)
}
