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

package global

import (
	"time"

	"github.com/cased/cased-go"
)

func NewPublisher(apiKey string) cased.Publisher {
	c := cased.NewPublisher(
		cased.WithTransport(cased.NewHTTPSyncTransport()),
		cased.WithSilence(false),
		cased.WithPublishKey(apiKey),

		// TODO: Replace with PublishURL once the API is ready
		// cased.WithPublishURL("https://api.meroxa.io/v1/telemetry"),
	)

	// The process will wait 30 seconds to publish all events to Cased before
	// exiting the process.
	defer c.Flush(30 * time.Second) // nolint:gomnd
	return c
}
