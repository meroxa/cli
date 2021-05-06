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

package open

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

var (
	_ builder.CommandWithDocs        = (*Open)(nil)
	_ builder.CommandWithSubCommands = (*Open)(nil)
)

type Open struct{}

func (o *Open) Usage() string {
	return "open"
}

func (o *Open) Docs() builder.Docs {
	return builder.Docs{
		Short: "Open in a web browser",
	}
}

func (o *Open) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Billing{}),
	}
}
