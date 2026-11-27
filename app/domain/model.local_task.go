package domain

// LocalTaskKinds are the kinds of [LocalTaskModel].
var LocalTaskKinds = []string{"enable", "disable", "remove", "purge", "local_ci"}

// States of a [LocalTaskModel].
const (
	// LocalTaskQueued: waiting to start.
	LocalTaskQueued = "queued"
	// LocalTaskRunning: under way.
	LocalTaskRunning = "running"
	// LocalTaskSucceeded: done.
	LocalTaskSucceeded = "succeeded"
	// LocalTaskFailed: given up; Error says why.
	LocalTaskFailed = "failed"
)

// LocalTaskModel is a job WarpEngine Local runs in the background:
// enabling, disabling or removing a provider, or setting a software up on
// the local CI.
type LocalTaskModel struct {
	// ID is the identifier.
	ID int `json:"id"`
	// Kind is one of [LocalTaskKinds].
	Kind string `json:"kind"`
	// Target is what it works on: a provider key or a software.
	Target string `json:"target"`
	// State is one of the LocalTask state constants
	// ([LocalTaskQueued] …).
	State string `json:"state"`
	// Steps are what it did so far, oldest first.
	Steps []LocalTaskStepModel `json:"steps"`
	// Result is what a finished task hands back, or nil; for a
	// local_ci: pushUrl, remoteName, repoUrl, ciRepositoryId, warnings.
	Result map[string]any `json:"result"`
	// Error is why it failed, or nil.
	Error *string `json:"error"`
	// CreatedAt is when it was queued (RFC 3339), or nil.
	CreatedAt *string `json:"createdAt"`
	// UpdatedAt is when it last changed (RFC 3339), or nil.
	UpdatedAt *string `json:"updatedAt"`
}

// Finished reports whether o has succeeded or failed.
func (o LocalTaskModel) Finished() bool {
	return o.State == LocalTaskSucceeded || o.State == LocalTaskFailed
}

// LocalTaskStepModel is one line of a [LocalTaskModel]'s progress.
type LocalTaskStepModel struct {
	// At is when it happened.
	At string `json:"at"`
	// Text says what happened.
	Text string `json:"text"`
}
