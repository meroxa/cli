package cased

import "time"

type Result struct {
	// The Result ID
	ID string `json:"id"`

	// The API URL for the result.
	APIURL string `json:"api_url"`

	// State contains the workflow run state.
	State WorkflowState `json:"state"`

	// Controls contains the controls specified by the workflow that was triggered.
	Controls ResultControls `json:"controls"`

	// Workflow contains the workflow that was detected or specified with the
	// event.
	Workflow *Workflow `json:"workflow"`

	// UpdatedAt is when the result was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CreatedAt is when the result was originally created.
	CreatedAt time.Time `json:"created_at"`
}

// ResultControls contains all the controls specified by the workflow that was
// triggered. All controls must be fulfilled for the workflow run to be
// fulfilled.
type ResultControls struct {
	Authentication *ResultControlsAuthentication `json:"authentication,omitempty"`
	Reason         *ResultControlsReason         `json:"reason,omitempty"`
	Approval       *ResultControlsApproval       `json:"approval,omitempty"`
}

type ResultControlsAuthentication struct {
	// State contains the authentication request state.
	State WorkflowState `json:"state"`

	// User contains the user information if the user has authenticated
	// successfully.
	User *ResultControlsAuthenticationUser `json:"user"`

	// URL contains the ephemeral URL for the user to authenticate with their
	// Cased account for this particular workflow run. If running in a headless
	// environment present the URL to the user to visit manually, otherwise
	// redirect the user to the URL.
	URL string `json:"url"`

	// ApiURL contains the URL to check the status of the authentication request.
	APIURL string `json:"api_url"`
}

type ResultControlsAuthenticationUser struct {
	// ID contains the authenticated Cased user ID.
	ID string `json:"id"`

	// Email contains the authenticated Cased user email address.
	Email string `json:"email"`
}

type ResultControlsReason struct {
	State WorkflowState `json:"state"`
}

// ResultControlsApprovalState reflects the workflow approval state.
type ResultControlsApprovalState string

const (
	// ResultControlsApprovalStatePending indicates when workflow controls have
	// not yet been fulfilled to request approval from the configured approval
	// sources.
	//
	// See that non-approval controls are fulfilled before the approval request
	// will be requested.
	ResultControlsApprovalStatePending ResultControlsApprovalState = "pending"

	// ResultControlsApprovalStateRequested indicates the workflow has requested
	// approval from the configured approval sources. The specified approval
	// requirements have not been met.
	ResultControlsApprovalStateRequested ResultControlsApprovalState = "requested"

	// ResultControlsApprovalStateApproved indicates that the approval request
	// has been approved per the configured approval requirements.
	ResultControlsApprovalStateApproved ResultControlsApprovalState = "approved"

	// ResultControlsApprovalStateDenied indicates that the approval request
	// has been denied or did not meet approval requirements.
	ResultControlsApprovalStateDenied ResultControlsApprovalState = "denied"

	// ResultControlsApprovalStateTimedOut indicates that the approval request
	// received no response within the configured timeout window.
	ResultControlsApprovalStateTimedOut ResultControlsApprovalState = "timed_out"

	// ResultControlsApprovalStateCanceled indicates that the approval request
	// was canceled by the requester.
	ResultControlsApprovalStateCanceled ResultControlsApprovalState = "canceled"
)

type ResultControlsApproval struct {
	// State reflects the workflow result's approval state.
	State ResultControlsApprovalState `json:"state"`

	// Requests contains all approval requests sent. Feature is not yet enabled.
	Requests []ResultControlsApprovalRequest `json:"requests"`

	// Source contains the approval sources.
	Source ResultControlsApprovalSource `json:"source"`
}

type ResultControlsApprovalRequestType string

type ResultControlsApprovalRequest struct {
	ID    string                            `json:"id"`
	State WorkflowState                     `json:"state"`
	Type  ResultControlsApprovalRequestType `json:"type"`
}

type ResultControlsApprovalSource struct {
	// Email indicates if the workflow approval request notified users via email.
	Email bool `json:"email"`

	// Slack indicates if the workflow approval request notified users via Slack.
	Slack ResultControlsApprovalSourceSlack `json:"slack"`
}

type ResultControlsApprovalSourceSlack struct {
	// Channel indicates the Slack channel the approval request was published to.
	Channel string `json:"channel"`
}
