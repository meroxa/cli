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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/global"
)

// needToCheckNewerCLIVersion checks different scenarios to determine whether to check or not
// 1. If user disabled auto-updating warning
// 2. If it checked within a week.
func needToCheckNewerCLIVersion() bool {
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

// getLatestCLIVersion returns latest CLI available tag.
func getLatestCLIVersion(ctx context.Context) (string, error) {
	client := &http.Client{}

	// Fetches tags in GitHub
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/repos/meroxa/cli/tags", http.NoBody)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result []struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return "", nil
	}

	return result[0].Name, nil
}
