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

package auth

import (
	"github.com/meroxa/cli/cmd/meroxa/builder"
	"github.com/spf13/cobra"
)

var (
	_ builder.CommandWithDocs        = (*Auth)(nil)
	_ builder.CommandWithSubCommands = (*Auth)(nil)
)

type Auth struct{}

func (o *Auth) Usage() string {
	return "auth"
}

func (o *Auth) Docs() builder.Docs {
	return builder.Docs{
		Short: "Authentication commands for Meroxa",
	}
}

func (o *Auth) SubCommands() []*cobra.Command {
	return []*cobra.Command{
		builder.BuildCobraCommand(&Login{}),
		builder.BuildCobraCommand(&Logout{}),
		builder.BuildCobraCommand(&WhoAmI{}),
	}
}
