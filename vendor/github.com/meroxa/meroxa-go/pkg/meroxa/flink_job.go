package meroxa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const flinkJobsBasePath = "/v1/flink-jobs"

type FlinkJobDesiredState string
type FlinkJobLifecycleState string
type FlinkJobReconciliationState string
type FlinkJobManagerDeploymentState string

const (
	FlinkJobLifecycleStateCreated       FlinkJobLifecycleState = "created"
	FlinkJobLifecycleStateDeploying     FlinkJobLifecycleState = "deploying"
	FlinkJobLifecycleStateDoa           FlinkJobLifecycleState = "doa"
	FlinkJobLifecycleStateFailed        FlinkJobLifecycleState = "failed"
	FlinkJobLifecycleStateRolledBack    FlinkJobLifecycleState = "rolled_back"
	FlinkJobLifecycleStateRollingBack   FlinkJobLifecycleState = "rolling_back"
	FlinkJobLifecycleStateStable        FlinkJobLifecycleState = "stable"
	FlinkJobLifecycleStateSuspended     FlinkJobLifecycleState = "suspended"
	FlinkJobLifecycleStateUninitialized FlinkJobLifecycleState = "uninitialized"
	FlinkJobLifecycleStateUpgrading     FlinkJobLifecycleState = "upgrading"

	FlinkJobDesiredStateRunning   FlinkJobDesiredState = "running"
	FlinkJobDesiredStateSuspended FlinkJobDesiredState = "suspended"

	FlinkJobReconciliationStateDeployed    FlinkJobReconciliationState = "deployed"
	FlinkJobReconciliationStateRolledBack  FlinkJobReconciliationState = "rolled_back"
	FlinkJobReconciliationStateRollingBack FlinkJobReconciliationState = "rolling_back"
	FlinkJobReconciliationStateUpgrading   FlinkJobReconciliationState = "upgrading"

	FlinkJobManagerDeploymentStateDeployedNotReady FlinkJobManagerDeploymentState = "deployed_not_ready"
	FlinkJobManagerDeploymentStateDeploying        FlinkJobManagerDeploymentState = "deploying"
	FlinkJobManagerDeploymentStateError            FlinkJobManagerDeploymentState = "error"
	FlinkJobManagerDeploymentStateFailing          FlinkJobManagerDeploymentState = "failing"
	FlinkJobManagerDeploymentStateMissing          FlinkJobManagerDeploymentState = "missing"
	FlinkJobManagerDeploymentStateReady            FlinkJobManagerDeploymentState = "ready"
)

type FlinkJobStatus struct {
	LifecycleState         FlinkJobLifecycleState         `json:"lifecycle_state"`
	State                  string                         `json:"state"`
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
	Name        string `json:"name"`
	JarURL      string `json:"jar_url"`
	Spec        string `json:"spec"`
	SpecVersion string `json:"spec_version"`
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
