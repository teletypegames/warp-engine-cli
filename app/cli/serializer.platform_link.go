package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// PlatformLinkSerializer turns the flags of the platform-link commands into
// the API's terms and the links into tables and records.
type PlatformLinkSerializer struct{}

// PlatformLinkSorts are the orders "platform-link list --sort" takes.
var PlatformLinkSorts = []string{"platform", "position", "name"}

// Filters are the filter flags of "platform-link list".
func (s PlatformLinkSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "platform", Key: "platform", Usage: "platform is"},
	}
}

// Fields are the field flags of "platform-link create" and "platform-link update".
func (s PlatformLinkSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "platform", Key: "platform", Usage: "platform, one of the server's"},
		{Flag: "name", Key: "name", Usage: "the link's text"},
		{Flag: "url", Key: "url", Usage: "where it points"},
		{Flag: "position", Key: "position", Kind: KindInt, Usage: "order within the platform, lowest first"},
	}
}

// FlagsToFilter reads the filter flags of "platform-link list".
func (s PlatformLinkSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of create and update.
func (s PlatformLinkSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders platform links as a table: ID, platform, position,
// name, URL.
func (s PlatformLinkSerializer) ModelToTable(links []domain.PlatformLinkModel) Table {
	table := Table{Headers: []string{"ID", "Platform", "Position", "Name", "URL"}}
	for _, link := range links {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(link.ID), link.Platform, strconv.Itoa(link.Position), Truncate(link.Name, 30), link.URL,
		})
	}
	return table
}

// ModelToRecord renders one platform link with every field.
func (s PlatformLinkSerializer) ModelToRecord(link domain.PlatformLinkModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(link.ID)},
		{Name: "Platform", Value: link.Platform},
		{Name: "Name", Value: link.Name},
		{Name: "URL", Value: link.URL},
		{Name: "Position", Value: strconv.Itoa(link.Position)},
		{Name: "Created", Value: Text(link.CreatedAt)},
		{Name: "Updated", Value: Text(link.UpdatedAt)},
	}
}
