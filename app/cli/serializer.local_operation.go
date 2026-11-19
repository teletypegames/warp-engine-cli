package cli

import (
	"sort"
	"strconv"
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// LocalOperationSerializer turns the background jobs of WarpEngine Local
// into records.
type LocalOperationSerializer struct{}

// ModelToRecord renders an operation: ID, kind, target, state, error, the
// times, one line per step and one per result key.
func (s LocalOperationSerializer) ModelToRecord(operation domain.LocalOperationModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(operation.ID)},
		{Name: "Kind", Value: operation.Kind},
		{Name: "Target", Value: operation.Target},
		{Name: "State", Value: operation.State},
		{Name: "Error", Value: Text(operation.Error)},
		{Name: "Created", Value: Text(operation.CreatedAt)},
		{Name: "Updated", Value: Text(operation.UpdatedAt)},
	}
	for _, step := range operation.Steps {
		record = append(record, Field{Name: "Step", Value: s.StepText(step)})
	}
	keys := make([]string, 0, len(operation.Result))
	for key := range operation.Result {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := operation.Result[key]
		text, ok := value.(string)
		if !ok {
			text = JSONText(value)
		}
		record = append(record, Field{Name: "Result " + key, Value: text})
	}
	return record
}

// StepText renders one step as "15:04:05  text": the time of an RFC 3339
// At in the local zone, or At as it is.
func (s LocalOperationSerializer) StepText(step domain.LocalOperationStepModel) string {
	at := step.At
	if parsed, err := time.Parse(time.RFC3339, step.At); err == nil {
		at = parsed.Local().Format("15:04:05")
	}
	return at + "  " + step.Text
}
