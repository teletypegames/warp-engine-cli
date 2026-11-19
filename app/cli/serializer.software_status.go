package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SoftwareStatusSerializer turns the flags of the status commands into the
// API's terms and the statuses into tables and records.
type SoftwareStatusSerializer struct{}

// SoftwareStatusSorts are the orders "status list --sort" takes.
var SoftwareStatusSorts = []string{"position", "key"}

// Filters are the filter flags of "status list".
func (s SoftwareStatusSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "default-visible", Key: "defaultVisible", Usage: "shown by the stores by default or not", Bool: true},
	}
}

// Fields are the field flags of "status create" and "status update".
func (s SoftwareStatusSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "key", Key: "key", Usage: "the word a software's status holds: a-z, 0-9, _ and -"},
		{Flag: "label", Key: "label", Usage: "display name"},
		{Flag: "position", Key: "position", Kind: KindInt, Usage: "order, lowest first"},
		{Flag: "default-visible", Key: "defaultVisible", Kind: KindBool, Usage: "the stores show such softwares by default"},
		{Flag: "badge", Key: "badge", Kind: KindNullable, Usage: "short tag the stores print on such a software"},
	}
}

// FlagsToFilter reads the filter flags of "status list".
func (s SoftwareStatusSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of create and update.
func (s SoftwareStatusSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders statuses as a table: ID, key, label, position,
// visible by default, badge, number of softwares.
func (s SoftwareStatusSerializer) ModelToTable(statuses []domain.SoftwareStatusModel) Table {
	table := Table{Headers: []string{"ID", "Key", "Label", "Position", "Visible", "Badge", "Softwares"}}
	for _, status := range statuses {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(status.ID), status.Key, Truncate(status.Label, 40), strconv.Itoa(status.Position),
			YesNo(status.DefaultVisible), Text(status.Badge), strconv.Itoa(status.SoftwareCount),
		})
	}
	return table
}

// ModelToRecord renders one status with every field.
func (s SoftwareStatusSerializer) ModelToRecord(status domain.SoftwareStatusModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(status.ID)},
		{Name: "Key", Value: status.Key},
		{Name: "Label", Value: status.Label},
		{Name: "Position", Value: strconv.Itoa(status.Position)},
		{Name: "Visible by default", Value: YesNo(status.DefaultVisible)},
		{Name: "Badge", Value: Text(status.Badge)},
		{Name: "Softwares", Value: strconv.Itoa(status.SoftwareCount)},
		{Name: "Created", Value: Text(status.CreatedAt)},
		{Name: "Updated", Value: Text(status.UpdatedAt)},
	}
}
