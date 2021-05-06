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

package auth

import (
	"context"
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/meroxa/cli/cmd/meroxa/global"
)

var (
	_ builder.CommandWithDocs    = (*Logout)(nil)
	_ builder.CommandWithExecute = (*Logout)(nil)
)

type Logout struct{}

func (l *Logout) Usage() string {
	return "logout"
}

func (l *Logout) Docs() builder.Docs {
	return builder.Docs{
		Short: "logout of the Meroxa platform",
	}
}

func (l *Logout) Execute(ctx context.Context) error {
	// TODO: add confirmation
	global.Config.Set("ACCESS_TOKEN", "")
	global.Config.Set("REFRESH_TOKEN", "")
	err := global.Config.WriteConfig()
	if err != nil {
		return err
	}
	fmt.Println("Successfully logged out.")
	return nil
}
