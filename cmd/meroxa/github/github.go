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

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var Client HTTPClient

// getContentHomebrewFormula returns from GitHub the content of the formula file for Meroxa CLI.
func getContentHomebrewFormula(ctx context.Context) (string, error) {
	url := "https://api.github.com/repos/meroxa/homebrew-taps/contents/Formula/meroxa.rb"
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	resp, err := Client.Do(req)
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

	type FormulaDefinition struct {
		Content []byte `json:"content,omitempty"`
	}

	var f FormulaDefinition
	if err := json.Unmarshal(b, &f); err != nil {
		return "", err
	}

	return bytes.NewBuffer(f.Content).String(), nil
}

// parseVersionFromFormulaFile receives the content of a formula file such as the one in
// "https://api.github.com/repos/meroxa/homebrew-taps/contents/Formula/meroxa.rb"
// and extracts its version number.
func parseVersionFromFormulaFile(content string) string {
	r := regexp.MustCompile(`version "(\d+.\d+.\d+)"`)
	matches := r.FindStringSubmatch(content)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// GetLatestCLITag fetches the content formula file from GitHub and then parses its version.
// example: 2.0.0
func GetLatestCLITag(ctx context.Context) (string, error) {
	brewFormulaFile, err := getContentHomebrewFormula(ctx)
	if err != nil {
		return brewFormulaFile, err
	}

	return parseVersionFromFormulaFile(brewFormulaFile), nil
}
