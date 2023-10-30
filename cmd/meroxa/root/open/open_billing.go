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

package open

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/log"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"

	"github.com/pkg/browser"
)

type Billing struct {
	logger log.Logger
}

var (
	_ builder.CommandWithDocs    = (*Billing)(nil)
	_ builder.CommandWithExecute = (*Billing)(nil)
	_ builder.CommandWithLogger  = (*Billing)(nil)
)

func (b *Billing) Usage() string {
	return "billing"
}

func (b *Billing) Docs() builder.Docs {
	return builder.Docs{
		Short: "Open your billing page in a web browser",
	}
}

func (b *Billing) Logger(logger log.Logger) {
	b.logger = logger
}

func (b *Billing) Execute(ctx context.Context) error {
	b.logger.Info(ctx, "Opening your billing page in your browser...")
	return browser.OpenURL(b.getBillingURL())
}

func (b *Billing) getBillingURL() string {
	return fmt.Sprintf("%s/settings/billing", global.GetMeroxaAPIURL())
}
