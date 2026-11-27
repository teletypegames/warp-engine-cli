package cli

import (
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SoftwareForkSerializer turns the flags of the fork commands into the API's
// terms and the announced forks into tables and records.
type SoftwareForkSerializer struct{}

// SoftwareForkSorts are the orders "fork list --sort" takes.
var SoftwareForkSorts = []string{"createdAt", "checkedAt"}

// Filters are the filter flags of "fork list". --software takes an ID or a
// name; the service turns it into softwareId.
func (s SoftwareForkSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "status", Key: "status", Usage: "status is (" + strings.Join(domain.SoftwareForkStatuses, ", ") + ")"},
		{Flag: "software", Key: "software", Usage: "parent software is (an ID or a name)"},
	}
}

// Fields are the field flags of "fork update".
func (s SoftwareForkSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "status", Key: "status", Usage: "visible lists the fork, hidden keeps it out", Enum: domain.SoftwareForkStatuses},
	}
}

// FlagsToFilter reads the filter flags of "fork list".
func (s SoftwareForkSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of "fork update".
func (s SoftwareForkSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders announced forks as a table: ID, status, parent software,
// engine, name, title, last check.
func (s SoftwareForkSerializer) ModelToTable(links []domain.SoftwareForkModel) Table {
	table := Table{Headers: []string{"ID", "Status", "Software", "Engine", "Name", "Title", "Checked"}}
	for _, link := range links {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(link.ID), link.Status, Text(link.SoftwareName), link.EngineURL, link.Name,
			Truncate(Text(link.Title), 40), Text(link.CheckedAt),
		})
	}
	return table
}

// ModelToRecord renders one announced fork with every field.
func (s SoftwareForkSerializer) ModelToRecord(link domain.SoftwareForkModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(link.ID)},
		{Name: "Status", Value: link.Status},
		{Name: "Software", Value: Text(link.SoftwareName)},
		{Name: "Software ID", Value: strconv.Itoa(link.SoftwareID)},
		{Name: "Engine", Value: link.EngineURL},
		{Name: "Name", Value: link.Name},
		{Name: "Title", Value: Text(link.Title)},
		{Name: "Checked", Value: Text(link.CheckedAt)},
		{Name: "Last error", Value: Text(link.LastError)},
		{Name: "Created", Value: Text(link.CreatedAt)},
		{Name: "Updated", Value: Text(link.UpdatedAt)},
	}
}
