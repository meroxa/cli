package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const flinkJobsBasePath = "/v1/flink-jobs"

type FlinkJobState string
type FlinkJobLifecycleState string
type FlinkJobReconciliationState string
type FlinkJobManagerDeploymentState string

const (
	FlinkJobLifecycleStateCreated       FlinkJobState = "created"
	FlinkJobLifecycleStateDeploying     FlinkJobState = "deploying"
	FlinkJobLifecycleStateDoa           FlinkJobState = "doa"
	FlinkJobLifecycleStateFailed        FlinkJobState = "failed"
	FlinkJobLifecycleStateRolledBack    FlinkJobState = "rolled_back"
	FlinkJobLifecycleStateRollingBack   FlinkJobState = "rolling_back"
	FlinkJobLifecycleStateStable        FlinkJobState = "stable"
	FlinkJobLifecycleStateSuspended     FlinkJobState = "suspended"
	FlinkJobLifecycleStateUninitialized FlinkJobState = "uninitialized"
	FlinkJobLifecycleStateUpgrading     FlinkJobState = "upgrading"

	FlinkJobStateRunning   FlinkJobState = "running"
	FlinkJobStateSuspended FlinkJobState = "suspended"

	FlinkJobReconciliationStateDeployed    FlinkJobState = "deployed"
	FlinkJobReconciliationStateRolledBack  FlinkJobState = "rolled_back"
	FlinkJobReconciliationStateRollingBack FlinkJobState = "rolling_back"
	FlinkJobReconciliationStateUpgrading   FlinkJobState = "upgrading"

	FlinkJobStateDeployedNotReady FlinkJobState = "deployed_not_ready"
	FlinkJobStateDeploying        FlinkJobState = "deploying"
	FlinkJobStateError            FlinkJobState = "error"
	FlinkJobStateFailing          FlinkJobState = "failing"
	FlinkJobStateMissing          FlinkJobState = "missing"
	FlinkJobStateReady            FlinkJobState = "ready"
)

type FlinkJobStatus struct {
	LifecycleState         FlinkJobLifecycleState         `json:"lifecycle_state"`
	State                  FlinkJobState                  `json:"state"`
	ReconciliationState    FlinkJobReconciliationState    `json:"reconciliation_state"`
	ManagerDeploymentState FlinkJobManagerDeploymentState `json:"manager_deployment_state"`
	Details                string                         `json:"details,omitempty"`
}

type FlinkJob struct {
	UUID          string           `json:"uuid"`
	Name          string           `json:"name"`
	InputStreams  []string         `json:"input_streams,omitempty"`
	OutputStreams []string         `json:"output_streams,omitempty"`
	Environment   EntityIdentifier `json:"environment,omitempty"`
	Status        FlinkJobStatus   `json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type CreateFlinkJobInput struct {
	Name   string `json:"name"`
	JarURL string `json:"jar_url"`
}

func (c *client) GetFlinkJob(ctx context.Context, nameOrUUID string) (*FlinkJob, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", flinkJobsBasePath, nameOrUUID), nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var fj *FlinkJob
	if err = json.NewDecoder(resp.Body).Decode(&fj); err != nil {
		return nil, err
	}

	return fj, nil
}

func (c *client) ListFlinkJobs(ctx context.Context) ([]*FlinkJob, error) {
	resp, err := c.MakeRequest(ctx, http.MethodGet, flinkJobsBasePath, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var fjs []*FlinkJob
	if err = json.NewDecoder(resp.Body).Decode(&fjs); err != nil {
		return nil, err
	}

	return fjs, nil
}

func (c *client) CreateFlinkJob(ctx context.Context, input *CreateFlinkJobInput) (*FlinkJob, error) {
	resp, err := c.MakeRequest(ctx, http.MethodPost, flinkJobsBasePath, input, nil, nil)
	if err != nil {
		return nil, err
	}

	if err = handleAPIErrors(resp); err != nil {
		return nil, err
	}

	var fj *FlinkJob
	if err = json.NewDecoder(resp.Body).Decode(&fj); err != nil {
		return nil, err
	}

	return fj, nil
}

func (c *client) DeleteFlinkJob(ctx context.Context, nameOrUUID string) error {
	resp, err := c.MakeRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/%s", flinkJobsBasePath, nameOrUUID), nil, nil, nil)
	if err != nil {
		return err
	}

	return handleAPIErrors(resp)
}
