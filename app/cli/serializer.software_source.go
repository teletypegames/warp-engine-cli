package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SoftwareSourceSerializer turns the flags of "source set" into the API's
// terms and a source into a record.
type SoftwareSourceSerializer struct{}

// Fields are the field flags of "source set".
func (s SoftwareSourceSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "url", Key: "url", Usage: "the git repository's address"},
		{Flag: "ref", Key: "ref", Kind: KindNullable, Usage: "branch or tag; empty for the default branch"},
	}
}

// RequestToAttributes reads the field flags of "source set".
func (s SoftwareSourceSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToRecord renders a source with every field.
func (s SoftwareSourceSerializer) ModelToRecord(source domain.SoftwareSourceModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(source.ID)},
		{Name: "Software", Value: strconv.Itoa(source.SoftwareID)},
		{Name: "URL", Value: source.URL},
		{Name: "Ref", Value: Text(source.Ref)},
		{Name: "Public", Value: Flag(source.Public)},
		{Name: "Verified", Value: Text(source.VerifiedAt)},
		{Name: "Verify error", Value: Text(source.VerifyError)},
		{Name: "Created", Value: Text(source.CreatedAt)},
		{Name: "Updated", Value: Text(source.UpdatedAt)},
	}
}
