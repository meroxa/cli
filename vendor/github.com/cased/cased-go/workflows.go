package cased

import "time"

type Workflow struct {
	// The Workflow ID
	ID string `json:"id"`

	// Name of the workflow to be used to trigger the workflow.
	Name *string `json:"name,omitempty"`

	// The API URL for the workflow.
	APIURL string `json:"api_url"`

	// Conditions are how Cased determines which workflow should run when an event
	// is published and a workflow is not specified.
	Conditions []WorkflowCondition `json:"conditions"`

	// Controls specifies the controls enabled for this workflow.
	Controls WorkflowControls `json:"controls"`

	// UpdatedAt is when the workflow was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt is when the workflow was created.
	CreatedAt time.Time `json:"created_at"`
}

type WorkflowState string

const (
	// Workflow
	WorkflowStatePending WorkflowState = "pending"

	// Workflow result state is unfulfilled when all controls have not been met.
	WorkflowStateUnfulfilled WorkflowState = "unfulfilled"

	// Workflow controls were not met. When the workflow result state is rejected
	// it's intended any further progress is canceled.
	WorkflowStateFulfilled WorkflowState = "fulfilled"

	// Workflow controls were not met. When the workflow result state is rejected
	// it's intended any further progress is canceled.
	WorkflowStateRejected WorkflowState = "rejected"
)

// WorkflowParams contains the available fields when creating and updating
// workflows.
type WorkflowParams struct {
	Params `json:"-"`

	// Name is optional and only required if you intend to trigger workflows
	// by publishing events directly to them.
	Name *string `json:"name,omitempty"`

	// Conditions specify the conditions the workflow should match when events
	// are not published directly to a workflow.
	Conditions []*WorkflowConditionParams `json:"conditions,omitempty"`

	// Configure the controls necessary for the workflow to reach the fulfilled
	// state.
	Controls *WorkflowControlsParams `json:"controls,omitempty"`
}

// WorkflowConditionOperator contains all condition operators available for
// workflows.
type WorkflowConditionOperator string

const (
	// WorkflowConditionOperatorEndsWith case-insensitive matches "world" in
	// "hello world"
	WorkflowConditionOperatorEndsWith WorkflowConditionOperator = "endsWith"

	// WorkflowConditionOperatorEqual case-insensitive matches "cased" both
	// "cased" or "Cased"
	WorkflowConditionOperatorEqual WorkflowConditionOperator = "eq"

	// WorkflowConditionOperatorIncludes case-insensitive matches when the value
	// is included in the specified field for both strings and arrays.
	WorkflowConditionOperatorIncludes WorkflowConditionOperator = "in"

	// WorkflowConditionOperatorNotEqual case-insensitive matches when the field
	// does not contain the specified value.
	WorkflowConditionOperatorNotEqual WorkflowConditionOperator = "not"

	// WorkflowConditionOperatorRegex matches based on the provided regular
	// expression. Not currently enabled.
	WorkflowConditionOperatorRegex WorkflowConditionOperator = "re"

	// WorkflowConditionOperatorStartsWith case-insensitive matches "hello" in
	// "hello world"
	WorkflowConditionOperatorStartsWith WorkflowConditionOperator = "startsWith"
)

// WorkflowCondition is an individual clause in one or more conditions that can be used
// to match incoming events.
//
// All conditions are evaluated ignoring the case of the value.
type WorkflowCondition struct {
	// The path to the field on the event to evaluate this condition for.
	Field string `json:"field"`

	// Operator specifies the operator use to evaluate the condition. See
	// `ConditionOperator` for all available operators.
	Operator WorkflowConditionOperator `json:"operator"`

	// Value contains the value to be used to evaluate the condition based on its
	// configured operator.
	Value string `json:"value"`
}

// WorkflowCondition is an individual clause in one or more conditions that can be used
// to match incoming events.
//
// All conditions are evaluated ignoring the case of the value.
type WorkflowConditionParams struct {
	// The path to the field on the event to evaluate this condition for.
	Field *string `json:"field"`

	// Operator specifies the operator use to evaluate the condition. See
	// `ConditionOperator` for all available operators.
	Operator *string `json:"operator"`

	// Value contains the value to be used to evaluate the condition based on its
	// configured operator.
	Value *string `json:"value"`
}

type WorkflowControls struct {
	// Require a user to provide a reason to continue the workflow.
	Reason *bool `json:"reason,omitempty"`

	// Require a user to authenticate with Cased to continue the workflow.
	Authentication *bool `json:"authentication,omitempty"`

	// Require a user to receive approval before a workflow is fulfilled or
	// rejected.
	Approval *WorkflowControlsApproval `json:"approval,omitempty"`
}

type WorkflowControlsParams struct {
	// Require a user to provide a reason to continue the workflow.
	Reason *bool `json:"reason,omitempty"`

	// Require a user to authenticate with Cased to continue the workflow.
	Authentication *bool `json:"authentication,omitempty"`

	// Require a user to receive approval before a workflow is fulfilled or
	// rejected.
	Approval *WorkflowControlsApprovalParams `json:"approval,omitempty"`
}

type WorkflowControlsApproval struct {
	// The number of approvals required to fulfill the approval requirement.
	//
	// Approval count cannot exceed the number of users on your account,
	// otherwise an error will be returned.
	Count int `json:"count"`

	// Permit an approval request to allow user requesting approval the ability
	// to approve their own request. If the Authentication control is disabled,
	// any user can approve the request and this setting is ignored.
	SelfApproval bool `json:"self_approval"`

	// Determine how long the approval lasts for.
	Duration int `json:"duration"`

	// Control how long the approval request is valid for. If not supplied,
	// approval requests can be responded to indefinitely.
	Timeout *int `json:"timeout"`

	// List of responders that can include individual users and groups of users
	// who are authorized to respond to the approval request.
	Responders *WorkflowControlsApprovalResponders `json:"responders,omitempty"`

	// Sources where to obtain the approval from. If not provided, defaults to
	// email.
	Sources *WorkflowControlsApprovalSources `json:"sources,omitempty"`
}

type WorkflowControlsApprovalParams struct {
	// The number of approvals required to fulfill the approval requirement.
	//
	// Approval count cannot exceed the number of users on your account,
	// otherwise an error will be returned.
	Count *int `json:"count,omitempty"`

	// Permit an approval request to allow user requesting approval the ability
	// to approve their own request. If the Authentication control is disabled,
	// any user can approve the request and this setting is ignored.
	SelfApproval *bool `json:"self_approval,omitempty"`

	// Determine how long the approval lasts for.
	Duration *int `json:"duration,omitempty"`

	// Control how long the approval request is valid for. If not supplied,
	// approval requests can be responded to indefinitely.
	Timeout *int `json:"timeout,omitempty"`

	// List of responders that can include individual users and groups of users
	// who are authorized to respond to the approval request.
	Responders *WorkflowControlsApprovalResponders `json:"responders,omitempty"`

	// Sources where to obtain the approval from. If not provided, defaults to
	// email.
	Sources *WorkflowControlsApprovalSourcesParams `json:"sources,omitempty"`
}

// WorkflowControlsApprovalResponders is the list of individual users and groups
// of users who are authorized to respond to an approval request.
type WorkflowControlsApprovalResponders map[string]string

// WorkflowControlsApprovalSources determines where approval requests are
// delivered.
type WorkflowControlsApprovalSources struct {
	// Email determines if an email is delivered for the approval request.
	Email bool `json:"email"`

	// Slack when provided, publishes a Slack message which users can respond to
	// the request.
	Slack *WorkflowControlsApprovalSourcesSlack `json:"slack,omitempty"`
}

// WorkflowControlsApprovalSourcesParams determines where approval requests are
// delivered.
type WorkflowControlsApprovalSourcesParams struct {
	// Email determines if an email is delivered for the approval request.
	Email *bool `json:"email,omitempty"`

	// Slack when provided, publishes a Slack message which users can respond to
	// the request.
	Slack *WorkflowControlsApprovalSourcesSlackParams `json:"slack,omitempty"`
}

// WorkflowControlsApprovalSourcesSlack configures which the Slack approval
// source.
type WorkflowControlsApprovalSourcesSlack struct {
	// The fully qualified Slack channel name (ie: #security).
	Channel string `json:"channel"`
}

// WorkflowControlsApprovalSourcesSlackParams configures which the Slack
// approval source.
type WorkflowControlsApprovalSourcesSlackParams struct {
	// The fully qualified Slack channel name (ie: #security).
	Channel *string `json:"channel,omitempty"`
}
