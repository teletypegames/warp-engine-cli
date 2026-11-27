package cli

import (
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// CiRepositorySerializer turns the flags of the CI repository commands into the
// API's terms and the CI repositories and their runs into tables and records.
type CiRepositorySerializer struct{}

// CiRepositorySorts are the orders "ci-repo list --sort" takes.
var CiRepositorySorts = []string{"repoName", "lastPipelineAt"}

// Filters are the filter flags of "ci-repo list".
func (s CiRepositorySerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "platform", Key: "platform", Usage: "platform is"},
		{Flag: "active", Key: "active", Usage: "still active on the CI server or not", Bool: true},
	}
}

// Fields are the field flags of "ci-repo update".
func (s CiRepositorySerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "platform", Key: "platform", Usage: "platform it builds for: one of the server's, or unknown"},
		{Flag: "software", Key: "softwareId", Kind: KindNullable, Usage: "the software (ID or name) its builds publish; one CI repository per software"},
	}
}

// FlagsToFilter reads the filter flags of "ci-repo list".
func (s CiRepositorySerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of "ci-repo update".
func (s CiRepositorySerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders CI repositories as a table: ID, repository, platform,
// active, software, last status, last run.
func (s CiRepositorySerializer) ModelToTable(ci_repositories []domain.CiRepositoryModel) Table {
	table := Table{Headers: []string{"ID", "Repository", "Platform", "Active", "Software", "Last status", "Last run"}}
	for _, ci_repository := range ci_repositories {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(ci_repository.ID), ci_repository.FullName, ci_repository.Platform, YesNo(ci_repository.Active),
			Number(ci_repository.SoftwareID), Text(ci_repository.LastPipelineStatus), Text(ci_repository.LastPipelineAt),
		})
	}
	return table
}

// ModelToRecord renders one CI repository with every field.
func (s CiRepositorySerializer) ModelToRecord(ci_repository domain.CiRepositoryModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(ci_repository.ID)},
		{Name: "Repository", Value: ci_repository.FullName},
		{Name: "Repository ID", Value: strconv.Itoa(ci_repository.RepoID)},
		{Name: "Platform", Value: ci_repository.Platform},
		{Name: "Active", Value: YesNo(ci_repository.Active)},
		{Name: "Software", Value: Number(ci_repository.SoftwareID)},
		{Name: "Last status", Value: Text(ci_repository.LastPipelineStatus)},
		{Name: "Last run", Value: Text(ci_repository.LastPipelineAt)},
		{Name: "Created", Value: Text(ci_repository.CreatedAt)},
		{Name: "Updated", Value: Text(ci_repository.UpdatedAt)},
	}
}

// SyncToTable renders the answer of a sync as one table: ID, change
// (created, updated, deactivated), repository, platform.
func (s CiRepositorySerializer) SyncToTable(sync domain.CiRepositorySyncModel) Table {
	table := Table{Headers: []string{"ID", "Change", "Repository", "Platform"}}
	groups := []struct {
		change          string
		ci_repositories []domain.CiRepositoryModel
	}{{"created", sync.Created}, {"updated", sync.Updated}, {"deactivated", sync.Deactivated}}
	for _, group := range groups {
		for _, ci_repository := range group.ci_repositories {
			table.Rows = append(table.Rows, []string{strconv.Itoa(ci_repository.ID), group.change, ci_repository.FullName, ci_repository.Platform})
		}
	}
	return table
}

// RunsToTable renders pipeline runs as a table: number, status, event,
// branch, commit (short), started, the first line of the message.
func (s CiRepositorySerializer) RunsToTable(runs []domain.PipelineRunModel) Table {
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
func (s CiRepositorySerializer) RunToRecord(run domain.PipelineRunModel) Record {
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
