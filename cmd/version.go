/*
Copyright © 2020 Meroxa Inc

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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

// VersionCmd represents the `meroxa version` command
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the Meroxa CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("meroxa/%s %s/%s\n", meroxaVersion, runtime.GOOS, runtime.GOARCH)
		},
	}
}

// Before changing this function, we'll need to update how the we're using the User-Agent when interacting with
// Platform-API: https://git.io/JtXCG
func VersionString() string {
	return fmt.Sprintf("Meroxa CLI %s", meroxaVersion)
}
