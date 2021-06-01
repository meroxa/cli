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
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/cased/cased-go"
)

func NewPublisher(options ...cased.PublisherOption) cased.Publisher {
	options = append(options, cased.WithTransport(cased.NewHTTPSyncTransport()))
	c := cased.NewPublisher(options...)

	// The process will wait 30 seconds to publish all events to Cased before
	// exiting the process.
	defer c.Flush(30 * time.Second) // nolint:gomnd
	return c
}

func addUserInfo(event *cased.AuditEvent) {
	actor, actorUUID, _ := GetCLIUserInfo()
	e := *event

	if actor != "" {
		e["actor"] = actor
	}

	if actorUUID != "" {
		e["actor_uuid"] = actorUUID
	}
}

func addAction(event *cased.AuditEvent, cmd *cobra.Command) {
	var action string
	e := *event

	// TODO: Implement something that could look up all the way up until meroxa (meroxa create resources...)
	// something like it determines how many levels since root and then until current cmd
	if cmd.HasParent() {
		if cmd.Parent().HasParent() {
			action = fmt.Sprintf("%s.%s.%s", cmd.Parent().Parent().Use, cmd.Parent().Use, cmd.Use)
		} else {
			action = fmt.Sprintf("%s.%s", cmd.Parent().Use, cmd.Use)
		}
	} else {
		action = cmd.Use
	}

	e["action"] = action
}

func BuildEvent(cmd *cobra.Command) cased.AuditEvent {
	event := cased.AuditEvent{
		"timestamp":  time.Now().UTC(),
		"user_agent": fmt.Sprintf("meroxa/%s %s/%s", Version, runtime.GOOS, runtime.GOARCH),
	}

	addUserInfo(&event)
	addAction(&event, cmd)

	return event
}
