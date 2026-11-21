package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// ExternalLinkSerializer turns the flags of the link commands into the
// API's terms and the links into tables and records.
type ExternalLinkSerializer struct{}

// Fields are the field flags of "link create" and "link update".
func (s ExternalLinkSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "label", Key: "label", Usage: "the link's text"},
		{Flag: "url", Key: "url", Usage: "where it points"},
		{Flag: "kind", Key: "kind", Usage: "source is the git repository, one per software", Enum: domain.ExternalLinkKinds},
		{Flag: "ref", Key: "ref", Kind: KindNullable, Usage: "branch or tag of a source; empty for the default branch"},
	}
}

// RequestToAttributes reads the field flags of create and update.
func (s ExternalLinkSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders links as a table: ID, kind, label, URL, ref,
// public.
func (s ExternalLinkSerializer) ModelToTable(links []domain.ExternalLinkModel) Table {
	table := Table{Headers: []string{"ID", "Kind", "Label", "URL", "Ref", "Public"}}
	for _, link := range links {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(link.ID), link.Kind, Truncate(link.Label, 40), link.URL, Text(link.Ref), Flag(link.Public),
		})
	}
	return table
}

// ModelToRecord renders one link with every field.
func (s ExternalLinkSerializer) ModelToRecord(link domain.ExternalLinkModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(link.ID)},
		{Name: "Software", Value: strconv.Itoa(link.SoftwareID)},
		{Name: "Label", Value: link.Label},
		{Name: "URL", Value: link.URL},
		{Name: "Kind", Value: link.Kind},
		{Name: "Ref", Value: Text(link.Ref)},
		{Name: "Public", Value: Flag(link.Public)},
		{Name: "Verified", Value: Text(link.VerifiedAt)},
		{Name: "Verify error", Value: Text(link.VerifyError)},
		{Name: "Created", Value: Text(link.CreatedAt)},
		{Name: "Updated", Value: Text(link.UpdatedAt)},
	}
}
