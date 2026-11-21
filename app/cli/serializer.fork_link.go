package cli

import (
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// ForkLinkSerializer turns the flags of the fork commands into the API's
// terms and the fork links into tables and records.
type ForkLinkSerializer struct{}

// ForkLinkSorts are the orders "fork list --sort" takes.
var ForkLinkSorts = []string{"createdAt", "checkedAt"}

// Filters are the filter flags of "fork list". --software takes an ID or a
// name; the service turns it into softwareId.
func (s ForkLinkSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "status", Key: "status", Usage: "status is (" + strings.Join(domain.ForkLinkStatuses, ", ") + ")"},
		{Flag: "software", Key: "software", Usage: "parent software is (an ID or a name)"},
	}
}

// Fields are the field flags of "fork update".
func (s ForkLinkSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "status", Key: "status", Usage: "visible lists the fork, hidden keeps it out", Enum: domain.ForkLinkStatuses},
	}
}

// FlagsToFilter reads the filter flags of "fork list".
func (s ForkLinkSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of "fork update".
func (s ForkLinkSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders fork links as a table: ID, status, parent software,
// child engine, child name, child title, last check.
func (s ForkLinkSerializer) ModelToTable(links []domain.ForkLinkModel) Table {
	table := Table{Headers: []string{"ID", "Status", "Software", "Engine", "Name", "Title", "Checked"}}
	for _, link := range links {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(link.ID), link.Status, Text(link.SoftwareName), link.ChildEngine, link.ChildName,
			Truncate(Text(link.ChildTitle), 40), Text(link.CheckedAt),
		})
	}
	return table
}

// ModelToRecord renders one fork link with every field.
func (s ForkLinkSerializer) ModelToRecord(link domain.ForkLinkModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(link.ID)},
		{Name: "Status", Value: link.Status},
		{Name: "Software", Value: Text(link.SoftwareName)},
		{Name: "Software ID", Value: strconv.Itoa(link.SoftwareID)},
		{Name: "Child engine", Value: link.ChildEngine},
		{Name: "Child name", Value: link.ChildName},
		{Name: "Child title", Value: Text(link.ChildTitle)},
		{Name: "Checked", Value: Text(link.CheckedAt)},
		{Name: "Last error", Value: Text(link.LastError)},
		{Name: "Created", Value: Text(link.CreatedAt)},
		{Name: "Updated", Value: Text(link.UpdatedAt)},
	}
}
