package cli

import (
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// PipelineSerializer turns the flags of the pipeline commands into the
// API's terms and the pipelines and their runs into tables and records.
type PipelineSerializer struct{}

// PipelineSorts are the orders "pipeline list --sort" takes.
var PipelineSorts = []string{"repoName", "lastPipelineAt"}

// Filters are the filter flags of "pipeline list".
func (s PipelineSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "platform", Key: "platform", Usage: "platform is"},
		{Flag: "active", Key: "active", Usage: "still active on the CI server or not", Bool: true},
	}
}

// Fields are the field flags of "pipeline update".
func (s PipelineSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "platform", Key: "platform", Usage: "platform it builds for", Enum: domain.PipelinePlatforms},
		{Flag: "software", Key: "softwareId", Kind: KindNullable, Usage: "the software (ID or name) its builds publish; one pipeline per software"},
	}
}

// FlagsToFilter reads the filter flags of "pipeline list".
func (s PipelineSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of "pipeline update".
func (s PipelineSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders pipelines as a table: ID, repository, platform,
// active, software, last status, last run.
func (s PipelineSerializer) ModelToTable(pipelines []domain.PipelineModel) Table {
	table := Table{Headers: []string{"ID", "Repository", "Platform", "Active", "Software", "Last status", "Last run"}}
	for _, pipeline := range pipelines {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(pipeline.ID), pipeline.FullName, pipeline.Platform, YesNo(pipeline.Active),
			Number(pipeline.SoftwareID), Text(pipeline.LastPipelineStatus), Text(pipeline.LastPipelineAt),
		})
	}
	return table
}

// ModelToRecord renders one pipeline with every field.
func (s PipelineSerializer) ModelToRecord(pipeline domain.PipelineModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(pipeline.ID)},
		{Name: "Repository", Value: pipeline.FullName},
		{Name: "Repository ID", Value: strconv.Itoa(pipeline.RepoID)},
		{Name: "Platform", Value: pipeline.Platform},
		{Name: "Active", Value: YesNo(pipeline.Active)},
		{Name: "Software", Value: Number(pipeline.SoftwareID)},
		{Name: "Last status", Value: Text(pipeline.LastPipelineStatus)},
		{Name: "Last run", Value: Text(pipeline.LastPipelineAt)},
		{Name: "Created", Value: Text(pipeline.CreatedAt)},
		{Name: "Updated", Value: Text(pipeline.UpdatedAt)},
	}
}

// SyncToTable renders the answer of a sync as one table: ID, change
// (created, updated, deactivated), repository, platform.
func (s PipelineSerializer) SyncToTable(sync domain.PipelineSyncModel) Table {
	table := Table{Headers: []string{"ID", "Change", "Repository", "Platform"}}
	groups := []struct {
		change    string
		pipelines []domain.PipelineModel
	}{{"created", sync.Created}, {"updated", sync.Updated}, {"deactivated", sync.Deactivated}}
	for _, group := range groups {
		for _, pipeline := range group.pipelines {
			table.Rows = append(table.Rows, []string{strconv.Itoa(pipeline.ID), group.change, pipeline.FullName, pipeline.Platform})
		}
	}
	return table
}

// RunsToTable renders pipeline runs as a table: number, status, event,
// branch, commit (short), started, the first line of the message.
func (s PipelineSerializer) RunsToTable(runs []domain.PipelineRunModel) Table {
	table := Table{Headers: []string{"Number", "Status", "Event", "Branch", "Commit", "Started", "Message"}}
	for _, run := range runs {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(run.Number), Text(run.Status), Text(run.Event), Text(run.Branch),
			shortCommit(run.Commit), Text(run.CreatedAt), Truncate(firstLine(run.Message), 50),
		})
	}
	return table
}

// RunToRecord renders one pipeline run with every field.
func (s PipelineSerializer) RunToRecord(run domain.PipelineRunModel) Record {
	return Record{
		{Name: "Number", Value: strconv.Itoa(run.Number)},
		{Name: "Status", Value: Text(run.Status)},
		{Name: "Event", Value: Text(run.Event)},
		{Name: "Branch", Value: Text(run.Branch)},
		{Name: "Ref", Value: Text(run.Ref)},
		{Name: "Commit", Value: Text(run.Commit)},
		{Name: "Started", Value: Text(run.CreatedAt)},
		{Name: "URL", Value: Text(run.URL)},
		{Name: "Message", Value: Text(run.Message)},
	}
}

// shortCommit renders a commit hash as its first 8 characters, or "-".
func shortCommit(commit *string) string {
	if commit == nil || *commit == "" {
		return "-"
	}
	if len(*commit) > 8 {
		return (*commit)[:8]
	}
	return *commit
}

// firstLine renders the first line of an optional text, trimmed, or "-".
func firstLine(text *string) string {
	if text == nil {
		return "-"
	}
	line, _, _ := strings.Cut(strings.TrimSpace(*text), "\n")
	if line == "" {
		return "-"
	}
	return strings.TrimSpace(line)
}
