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
package main

import (
	"fmt"

	"github.com/meroxa/cli/cmd/meroxa/global"
	"github.com/meroxa/cli/cmd/meroxa/root"
)

var (
	version      = "dev"
	GitCommit    string
	GitUntracked string
	GitLatestTag string
)

func gitInfoNotEmpty() bool {
	return GitCommit != "" || GitLatestTag != "" || GitUntracked != ""
}

func main() {
	// In production this will be the updated version by `goreleaser`
	// We need to include `v` since we'll compare with the actual tag in GitHub
	// For more information see https://goreleaser.com/cookbooks/using-main.version
	global.CurrentTag = fmt.Sprintf("v%s", version)

	if gitInfoNotEmpty() {
		if GitCommit != "" {
			version += fmt.Sprintf(":%s", GitCommit)
		}

		if GitLatestTag != "" {
			version += fmt.Sprintf(" %s", GitLatestTag)
		}

		if GitUntracked != "" {
			version += fmt.Sprintf(" %s", GitUntracked)
		}
	}

	global.Version = version
	root.Run()
}
