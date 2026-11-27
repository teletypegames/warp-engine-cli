package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SoftwareLinkSerializer turns the flags of the link commands into the
// API's terms and the links into tables and records.
type SoftwareLinkSerializer struct{}

// Fields are the field flags of "link create" and "link update".
func (s SoftwareLinkSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "label", Key: "label", Usage: "the link's text"},
		{Flag: "url", Key: "url", Usage: "where it points"},
	}
}

// RequestToAttributes reads the field flags of create and update.
func (s SoftwareLinkSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders links as a table: ID, label, URL.
func (s SoftwareLinkSerializer) ModelToTable(links []domain.SoftwareLinkModel) Table {
	table := Table{Headers: []string{"ID", "Label", "URL"}}
	for _, link := range links {
		table.Rows = append(table.Rows, []string{strconv.Itoa(link.ID), Truncate(link.Label, 40), link.URL})
	}
	return table
}

// ModelToRecord renders one link with every field.
func (s SoftwareLinkSerializer) ModelToRecord(link domain.SoftwareLinkModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(link.ID)},
		{Name: "Software", Value: strconv.Itoa(link.SoftwareID)},
		{Name: "Label", Value: link.Label},
		{Name: "URL", Value: link.URL},
		{Name: "Created", Value: Text(link.CreatedAt)},
		{Name: "Updated", Value: Text(link.UpdatedAt)},
	}
}
