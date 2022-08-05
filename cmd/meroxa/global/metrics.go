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

	casedAPIKey := Config.GetString(CasedPublishKeyEnv)

	if casedAPIKey != "" {
		options = append(options, cased.WithPublishKey(casedAPIKey))
	} else {
		options = append(options, cased.WithPublishURL(fmt.Sprintf("%s/telemetry", GetMeroxaAPIURL())))
	}

	if v := Config.GetString(PublishMetricsEnv); v == "false" {
		options = append(options, cased.WithSilence(true))
	}

	if v := Config.GetBool(CasedDebugEnv); v {
		options = append(options, cased.WithDebug(v))
	}

	// TODO: Re-use HTTP Transform from our own Meroxa client
	options = append(options, cased.WithTransport(cased.NewHTTPTransport()))

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

func addError(event cased.AuditEvent, err error) {
	if err != nil {
		event["error"] = err.Error()
	}
}

func addArgs(args []string) string {
	if len(args) > 0 {
		return strings.Join(args, ",")
	}

	return ""
}

func addUserInfo(event cased.AuditEvent) {
	actor, actorUUID, _ := GetCLIUserInfo()

	if actor != "" {
		event["actor"] = actor
	}

	if actorUUID != "" {
		event["actor_uuid"] = actorUUID
	}
}

// addAction looks up all parent commands recursively and builds the action field.
func addAction(event cased.AuditEvent, cmd *cobra.Command) {
	var buildAction func(cmd *cobra.Command) string
	buildAction = func(cmd *cobra.Command) string {
		if !cmd.HasParent() {
			return cmd.Use
		}
		action := buildAction(cmd.Parent())
		return fmt.Sprintf("%s.%s", action, cmd.Use)
	}

	event["action"] = buildAction(cmd)
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

	addUserInfo(event)
	addAction(event, cmd)
	addError(event, err)

	return event
}

var (
	// PublishEvent is a public variable so in builder tests we could overwrite it.
	PublishEvent = publishEvent
)

// PublishEvent will take care of publishing the event to Cased.
func publishEvent(event cased.AuditEvent) {
	// Only prints out to console
	if v := Config.GetString(PublishMetricsEnv); v == "stdout" {
		e, _ := json.Marshal(event)
		fmt.Printf("\n\nEvent: %v\n\n", string(e))
		return
	}

	publisher := NewPublisher()
	cased.SetPublisher(publisher)

	err := cased.Publish(event)
	if err != nil {
		// cased.Publish could return an error, but we only show it when debugging
		if Config.GetBool(CasedDebugEnv) {
			fmt.Println("error: %w", err)
		}
		return
	}

	// The process will wait 30 seconds to publish all events to Cased before
	// exiting the process.
	cased.Flush(30 * time.Second) //nolint:gomnd
}
