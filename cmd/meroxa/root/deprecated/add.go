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

/* ⚠️ WARN ⚠️

The following commands will be removed once we decide to stop adding support for commands that don't follow
the `subject-verb-object` design.

*/

package deprecated

import (
	"github.com/meroxa/cli/cmd/meroxa/global"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/root/resources"
	"github.com/spf13/cobra"
)

// addCmd represents `meroxa add`.
func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a resource to your Meroxa resource catalog",
	}

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `[connector | endpoint | pipeline | resource] create` instead"
	}

	cmd.AddCommand(addResourceCmd())
	return cmd
}

// addResourceCmd represents `meroxa add resource` -> `meroxa resources create`.
func addResourceCmd() *cobra.Command {
	cmd := builder.BuildCobraCommand(&resources.Create{})
	cmd.Use = "resource [NAME] --type TYPE --url URL"
	cmd.Short = "Add a resource to your Meroxa resource catalog"
	cmd.Long = `Use the add command to add resources to your Meroxa resource catalog.`
	cmd.Example = `
meroxa add resource store --type postgres -u $DATABASE_URL --metadata '{"logical_replication":true}'
meroxa add resource datalake --type s3 -u "s3://$AWS_ACCESS_KEY_ID:$AWS_ACCESS_KEY_SECRET@us-east-1/meroxa-demos"
meroxa add resource warehouse --type redshift -u $REDSHIFT_URL
meroxa add resource slack --type url -u $WEBHOOK_URL
`
	cmd.Aliases = []string{"resources"}

	if global.IsMeroxaV2Released() {
		cmd.Deprecated = "use `resource create` instead"
	}

	return cmd
}
