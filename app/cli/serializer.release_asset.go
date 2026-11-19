package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// ReleaseAssetSerializer turns the flags of the asset commands into the
// API's terms and the assets into tables and records.
type ReleaseAssetSerializer struct{}

// Fields are the field flags of "asset create" and "asset update".
func (s ReleaseAssetSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "kind", Key: "kind", Usage: "what the file is", Enum: domain.AssetKinds},
		{Flag: "path", Key: "path", Usage: "absolute path of the file in the artifact store"},
	}
}

// RequestToAttributes reads the field flags of create and update.
func (s ReleaseAssetSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders assets as a table: ID, release, kind, path.
func (s ReleaseAssetSerializer) ModelToTable(assets []domain.ReleaseAssetModel) Table {
	table := Table{Headers: []string{"ID", "Release", "Kind", "Path"}}
	for _, asset := range assets {
		table.Rows = append(table.Rows, []string{strconv.Itoa(asset.ID), strconv.Itoa(asset.ReleaseID), asset.Kind, asset.Path})
	}
	return table
}

// ModelToRecord renders one asset with every field.
func (s ReleaseAssetSerializer) ModelToRecord(asset domain.ReleaseAssetModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(asset.ID)},
		{Name: "Release", Value: strconv.Itoa(asset.ReleaseID)},
		{Name: "Kind", Value: asset.Kind},
		{Name: "Path", Value: asset.Path},
		{Name: "Created", Value: Text(asset.CreatedAt)},
		{Name: "Updated", Value: Text(asset.UpdatedAt)},
	}
}
