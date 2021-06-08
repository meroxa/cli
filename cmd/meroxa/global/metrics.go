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
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/cased/cased-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewPublisher() cased.Publisher {
	var options []cased.PublisherOption

	casedAPIKey := Config.GetString("CASED_API_KEY")

	if casedAPIKey != "" {
		options = append(options, cased.WithPublishKey(casedAPIKey))
	} else {
		options = append(options, cased.WithPublishURL(fmt.Sprintf("%s/telemetry", GetMeroxaAPIURL())))
	}

	if v := Config.GetString("PUBLISH_METRICS"); v == "false" {
		options = append(options, cased.WithSilence(true))
	}

	if v := Config.GetBool("CASED_DEBUG"); v {
		options = append(options, cased.WithDebug(v))
	}
	c := cased.NewPublisher(options...)
	return c
}

func addAlias(cmd *cobra.Command) string {
	if cmd.Use != cmd.CalledAs() {
		return cmd.CalledAs()
	}

	return ""
}

func addFlags(cmd *cobra.Command) string {
	var flags []string

	if cmd.HasFlags() {
		cmd.Flags().Visit(func(flag *pflag.Flag) {
			flags = append(flags, flag.Name)
		})
	}

	return strings.Join(flags, ",")
}

func addDeprecated(cmd *cobra.Command) bool {
	return cmd.Deprecated != ""
}

func addError(event *cased.AuditEvent, err error) {
	e := *event

	if err != nil {
		e["error"] = err.Error()
	}
}

func addArgs(args []string) string {
	if len(args) > 0 {
		return strings.Join(args, ",")
	}

	return ""
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

func buildCommandInfo(cmd *cobra.Command, args []string) map[string]interface{} {
	c := make(map[string]interface{})

	if v := addAlias(cmd); v != "" {
		c["alias"] = v
	}

	if v := addArgs(args); v != "" {
		c["args"] = v
	}

	if v := addFlags(cmd); v != "" {
		c["flags"] = v
	}

	if v := addDeprecated(cmd); v {
		c["deprecated"] = v
	}

	return c
}

func buildBasicEvent(cmd *cobra.Command, args []string) cased.AuditEvent {
	return cased.AuditEvent{
		"timestamp":  time.Now().UTC(),
		"user_agent": fmt.Sprintf("meroxa/%s %s/%s", Version, runtime.GOOS, runtime.GOARCH),
		"command":    buildCommandInfo(cmd, args),
	}
}

func BuildEvent(cmd *cobra.Command, args []string, err error) cased.AuditEvent {
	event := buildBasicEvent(cmd, args)

	addUserInfo(&event)
	addAction(&event, cmd)
	addError(&event, err)

	return event
}

var (
	// PublishEvent is a public variable so in builder tests we could overwrite it.
	PublishEvent = publishEvent
)

// PublishEvent will take care of publishing the event to Cased.
func publishEvent(event cased.AuditEvent) {
	// Only prints out to console
	if v := Config.GetString("PUBLISH_METRICS"); v == "stdout" {
		e, _ := json.Marshal(event)
		fmt.Printf("\n\nEvent: %v\n\n", string(e))
		return
	}

	publisher := NewPublisher()
	cased.SetPublisher(publisher)

	// cased.Publish could return an error, but we only show it when debugging
	err := cased.Publish(event)

	// The process will wait 30 seconds to publish all events to Cased before
	// exiting the process.
	defer cased.Flush(30 * time.Second) // nolint:gomnd

	if err != nil && Config.GetBool("CASED_DEBUG") {
		fmt.Println("error: %w", err)
	}
}
