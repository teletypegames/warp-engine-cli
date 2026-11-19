package domain

// LocalOperationKinds are the kinds of [LocalOperationModel].
var LocalOperationKinds = []string{"enable", "disable", "remove", "purge", "local_ci"}

// States of a [LocalOperationModel].
const (
	// LocalOperationQueued: waiting to start.
	LocalOperationQueued = "queued"
	// LocalOperationRunning: under way.
	LocalOperationRunning = "running"
	// LocalOperationSucceeded: done.
	LocalOperationSucceeded = "succeeded"
	// LocalOperationFailed: given up; Error says why.
	LocalOperationFailed = "failed"
)

// LocalOperationModel is a job WarpEngine Local runs in the background:
// enabling, disabling or removing a component, or setting a software up on
// the local CI.
type LocalOperationModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Kind is one of [LocalOperationKinds].
	Kind string `json:"kind"`
	// Target is what it works on: a component key or a software.
	Target string `json:"target"`
	// State is one of the LocalOperation state constants
	// ([LocalOperationQueued] …).
	State string `json:"state"`
	// Steps are what it did so far, oldest first.
	Steps []LocalOperationStepModel `json:"steps"`
	// Result is what a finished operation hands back, or nil; for a
	// local_ci: pushUrl, remoteName, repoUrl, pipelineId, warnings.
	Result map[string]any `json:"result"`
	// Error is why it failed, or nil.
	Error *string `json:"error"`
	// CreatedAt is when it was queued (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}

// Finished reports whether o has succeeded or failed.
func (o LocalOperationModel) Finished() bool {
	return o.State == LocalOperationSucceeded || o.State == LocalOperationFailed
}

// LocalOperationStepModel is one line of a [LocalOperationModel]'s progress.
type LocalOperationStepModel struct {
	// At is when it happened.
	At string `json:"at"`
	// Text says what happened.
	Text string `json:"text"`
}
