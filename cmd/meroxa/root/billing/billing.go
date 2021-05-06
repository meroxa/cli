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

package billing

import (
	"context"

	"github.com/meroxa/cli/cmd/meroxa/builder"

	"github.com/meroxa/cli/cmd/meroxa/root/open"
)

type Billing struct{}

var (
	_ builder.CommandWithDocs    = (*Billing)(nil)
	_ builder.CommandWithExecute = (*Billing)(nil)
)

func (b *Billing) Usage() string {
	return "billing"
}

func (b *Billing) Docs() builder.Docs {
	return builder.Docs{
		Short: "Open your billing page in a web browser",
	}
}

func (b *Billing) Execute(ctx context.Context) error {
	err := builder.BuildCobraCommand(&open.Billing{}).Execute()

	if err != nil {
		return err
	}

	return nil
}
