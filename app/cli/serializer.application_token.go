package cli

import (
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// ApplicationTokenSerializer turns the flags of the token commands into the
// API's terms and the application tokens into tables and records.
type ApplicationTokenSerializer struct{}

// ApplicationTokenSorts are the orders "token list --sort" takes.
var ApplicationTokenSorts = []string{"createdAt", "lastUsedAt", "name"}

// Filters are the filter flags of "token list".
func (s ApplicationTokenSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "state", Key: "state", Usage: "state is (" + strings.Join(domain.ApplicationTokenStates, ", ") + ")"},
		{Flag: "kind", Key: "kind", Usage: "kind is (" + strings.Join(domain.ApplicationTokenKinds, ", ") + ")"},
	}
}

// Fields are the field flags of "token update".
func (s ApplicationTokenSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "name", Key: "name", Usage: "what uses it"},
		{Flag: "scopes", Key: "scopes", Kind: KindStrings, Usage: "scopes, no wider than your own (admin:read, admin:write, admin:tokens …)"},
		{Flag: "unrestricted", Key: "unrestricted", Kind: KindBool, Usage: "reach every owner's records (only an unrestricted caller may grant it)"},
		{Flag: "expires-at", Key: "expiresAt", Kind: KindNullable, Usage: "expiry as RFC 3339 time, not later than your own"},
	}
}

// CreateFields are the field flags of "token create": the fields of
// update and --owner-id. Without --expires-at the token expires with yours.
func (s ApplicationTokenSerializer) CreateFields() []FieldFlag {
	return append(s.Fields(), FieldFlag{Flag: "owner-id", Key: "ownerId", Kind: KindInt, Usage: "the owning account (default: yours)"})
}

// FlagsToFilter reads the filter flags of "token list".
func (s ApplicationTokenSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// CreateToAttributes reads the field flags of create.
func (s ApplicationTokenSerializer) CreateToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.CreateFields())
}

// RequestToAttributes reads the field flags of update.
func (s ApplicationTokenSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders tokens as a table: ID, name, prefix, scopes,
// unrestricted, expiry, last use.
func (s ApplicationTokenSerializer) ModelToTable(tokens []domain.ApplicationTokenModel) Table {
	table := Table{Headers: []string{"ID", "Name", "Prefix", "Scopes", "Unrestricted", "Expires", "Last used"}}
	for _, token := range tokens {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(token.ID), token.Name, token.TokenPrefix, List(token.Scopes), YesNo(token.Unrestricted),
			Text(token.ExpiresAt), Text(token.LastUsedAt),
		})
	}
	return table
}

// ModelToRecord renders one token with every field but the plain token,
// which is printed apart from it; after a create or a rotate, the CI
// secrets it was written to.
func (s ApplicationTokenSerializer) ModelToRecord(token domain.ApplicationTokenModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(token.ID)},
		{Name: "Name", Value: token.Name},
		{Name: "Prefix", Value: token.TokenPrefix},
		{Name: "Owner", Value: token.OwnerType + " #" + strconv.Itoa(token.OwnerID)},
		{Name: "Scopes", Value: List(token.Scopes)},
		{Name: "Unrestricted", Value: YesNo(token.Unrestricted)},
		{Name: "Expires", Value: Text(token.ExpiresAt)},
		{Name: "Last used", Value: Text(token.LastUsedAt)},
		{Name: "Created", Value: Text(token.CreatedAt)},
		{Name: "Updated", Value: Text(token.UpdatedAt)},
	}
	if token.Sync != nil {
		record = append(record, Field{Name: "CI secrets", Value: strconv.Itoa(token.Sync.Synced) + " updated, " + strconv.Itoa(token.Sync.Failed) + " failed"})
	}
	return record
}
