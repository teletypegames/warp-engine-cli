package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SyncPairingSerializer turns the flags of the pairing commands into the
// API's terms and the pairings into tables and records.
type SyncPairingSerializer struct{}

// Fields are the field flags of "pairing create" and "pairing update".
func (s SyncPairingSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "label", Key: "label", Usage: "what the other engine is, e.g. \"laptop local\""},
	}
}

// RequestToAttributes reads the field flags of create and update.
func (s SyncPairingSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders pairings as a table: ID, label, token prefix, last
// inbound push, created.
func (s SyncPairingSerializer) ModelToTable(pairings []domain.SyncPairingModel) Table {
	table := Table{Headers: []string{"ID", "Label", "Token", "Last inbound", "Created"}}
	for _, pairing := range pairings {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(pairing.ID), pairing.Label, Text(pairing.InboundTokenPrefix), Text(pairing.LastInboundAt), Text(pairing.CreatedAt),
		})
	}
	return table
}

// ModelToRecord renders one pairing; its plain token, when the answer has
// one, is printed apart from it.
func (s SyncPairingSerializer) ModelToRecord(pairing domain.SyncPairingModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(pairing.ID)},
		{Name: "Label", Value: pairing.Label},
		{Name: "Token prefix", Value: Text(pairing.InboundTokenPrefix)},
		{Name: "Last inbound", Value: Text(pairing.LastInboundAt)},
		{Name: "Created", Value: Text(pairing.CreatedAt)},
		{Name: "Updated", Value: Text(pairing.UpdatedAt)},
	}
}
