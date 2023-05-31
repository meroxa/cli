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

package flink

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/meroxa/cli/utils/display"

	"github.com/golang/mock/gomock"
	"github.com/meroxa/cli/log"
	"github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go/pkg/meroxa"
	"github.com/meroxa/meroxa-go/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func getFlinkJobs() []*meroxa.FlinkJob {
	var flinkJobs []*meroxa.FlinkJob
	f := utils.GenerateFlinkJob()
	return append(flinkJobs, &f)
}

func TestListFlinkJobsExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	flinkJobs := append(getFlinkJobs())

	client.
		EXPECT().
		ListFlinkJobs(ctx).
		Return(flinkJobs, nil)

	l := &List{
		client: client,
		logger: logger,
	}

	err := l.Execute(ctx)
	if err != nil {
		t.Fatalf("not expected error, got \"%s\"", err.Error())
	}

	gotLeveledOutput := logger.LeveledOutput()
	wantLeveledOutput := display.FlinkJobsTable(flinkJobs)

	if !strings.Contains(gotLeveledOutput, wantLeveledOutput) {
		t.Fatalf("expected output:\n%s\ngot:\n%s", wantLeveledOutput, gotLeveledOutput)
	}

	gotJSONOutput := logger.JSONOutput()
	var gotJobs []meroxa.FlinkJob
	err = json.Unmarshal([]byte(gotJSONOutput), &gotJobs)

	var lf []meroxa.FlinkJob

	for _, f := range flinkJobs {
		lf = append(lf, *f)
	}

	if err != nil {
		t.Fatalf("not expected error, got %q", err.Error())
	}

	if !reflect.DeepEqual(gotJobs, lf) {
		t.Fatalf("expected \"%v\", got \"%v\"", lf, gotJobs)
	}
}
func TestListFlinkJobsErrorHandling(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	logger := log.NewTestLogger()

	errMsg := "some API error"

	client.
		EXPECT().
		ListFlinkJobs(ctx).
		Return(nil, errors.New(errMsg))

	l := &List{
		client: client,
		logger: logger,
	}

	err := l.Execute(ctx)
	if err == nil {
		t.Fatalf("expected error, got %q", err.Error())
	}

	assert.ErrorContains(t, err, errMsg)
}
