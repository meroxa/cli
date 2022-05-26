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

package builder

import (
	"time"

	"github.com/meroxa/cli/cmd/meroxa/global"
)

// needToCheckNewerCLIVersion checks different scenarios to determine whether to check or not
// 1. If user disabled auto-updating warning
// 2. If it checked within a week.
func needToCheckNewerCLIVersion() bool {
	if global.Config == nil {
		return false
	}

	disabledNotificationsUpdate := global.Config.GetBool(global.DisableNotificationsUpdate)
	if disabledNotificationsUpdate {
		return false
	}

	latestUpdatedAt := global.Config.GetTime(global.LatestCLIVersionUpdatedAtEnv)
	if latestUpdatedAt.IsZero() {
		return true
	}

	duration := time.Now().UTC().Sub(latestUpdatedAt)
	return duration.Hours() > 24*7
}

// getCurrentCLIVersion returns current CLI tag (example: `v2.0.0`)
// version, set by GoReleaser + `v` at the beginning.
func getCurrentCLIVersion() string {
	return global.CurrentTag
}
