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
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/meroxa/cli/cmd/meroxa/global"
)

// needToCheckNewerCLIVersion checks if CLI within a week
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
	return duration.Hours() > 24*7 // nolint:gomnd
}

// getCurrentCLIVersion returns current CLI tag
func getCurrentCLIVersion() string {
	return global.CurrentTag
}

type TagResponse struct {
	Name string `json:"name"`
}

// getLatestCLIVersion returns latest CLI available tag
func getLatestCLIVersion() string {
	client := &http.Client{}

	// Fetches tags in GitHub
	req, err := http.NewRequest("GET", "https://api.github.com/repos/meroxa/cli/tags", nil)
	if err != nil {
		return ""
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result []TagResponse
	if err := json.Unmarshal(b, &result); err != nil {
		return ""
	}

	return result[0].Name
}
