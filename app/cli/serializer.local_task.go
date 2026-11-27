package cli

import (
	"sort"
	"strconv"
	"time"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// LocalTaskSerializer turns the background jobs of WarpEngine Local
// into records.
type LocalTaskSerializer struct{}

// ModelToRecord renders a task: ID, kind, target, state, error, the
// times, one line per step and one per result key.
func (s LocalTaskSerializer) ModelToRecord(task domain.LocalTaskModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(task.ID)},
		{Name: "Kind", Value: task.Kind},
		{Name: "Target", Value: task.Target},
		{Name: "State", Value: task.State},
		{Name: "Error", Value: Text(task.Error)},
		{Name: "Created", Value: Text(task.CreatedAt)},
		{Name: "Updated", Value: Text(task.UpdatedAt)},
	}
	for _, step := range task.Steps {
		record = append(record, Field{Name: "Step", Value: s.StepText(step)})
	}
	keys := make([]string, 0, len(task.Result))
	for key := range task.Result {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := task.Result[key]
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
func (s LocalTaskSerializer) StepText(step domain.LocalTaskStepModel) string {
	at := step.At
	if parsed, err := time.Parse(time.RFC3339, step.At); err == nil {
		at = parsed.Local().Format("15:04:05")
	}
	return at + "  " + step.Text
}
