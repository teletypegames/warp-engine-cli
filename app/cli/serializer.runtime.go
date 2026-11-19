package cli

import (
	"sort"
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// RuntimeSerializer turns the flags of the runtime commands into the API's
// terms and the runtimes into tables and records.
type RuntimeSerializer struct{}

// RuntimeSorts are the orders "runtime list --sort" takes.
var RuntimeSorts = []string{"name", "createdAt"}

// Filters are the filter flags of "runtime list".
func (s RuntimeSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "name", Key: "name", Usage: "runtime name is"},
	}
}

// FlagsToFilter reads the filter flags of "runtime list".
func (s RuntimeSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// ModelToTable renders runtime versions as a table: ID, name, version,
// source, size, releases bound, installed.
func (s RuntimeSerializer) ModelToTable(runtimes []domain.RuntimeVersionModel) Table {
	table := Table{Headers: []string{"ID", "Name", "Version", "Source", "Size", "Releases", "Installed"}}
	for _, runtime := range runtimes {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(runtime.ID), runtime.Name, runtime.Version, runtime.Source,
			Bytes(runtime.SizeBytes), strconv.Itoa(runtime.Releases), Text(runtime.CreatedAt),
		})
	}
	return table
}

// ModelToRecord renders one runtime version with every field, one line per
// component.
func (s RuntimeSerializer) ModelToRecord(runtime domain.RuntimeVersionModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(runtime.ID)},
		{Name: "Name", Value: runtime.Name},
		{Name: "Version", Value: runtime.Version},
		{Name: "Source", Value: runtime.Source},
		{Name: "Path", Value: runtime.Path},
		{Name: "Size", Value: Bytes(runtime.SizeBytes)},
		{Name: "Releases", Value: strconv.Itoa(runtime.Releases)},
		{Name: "Installed", Value: Text(runtime.CreatedAt)},
		{Name: "Updated", Value: Text(runtime.UpdatedAt)},
	}
	for _, component := range runtime.Components {
		record = append(record, Field{Name: "Component", Value: component.Kind + " " + component.Key})
	}
	return record
}

// InstallToRecord renders the answer of an install: the runtime, whether
// it was a dry run, the files fetched and kept, the components by kind.
func (s RuntimeSerializer) InstallToRecord(install domain.RuntimeInstallModel) Record {
	record := Record{
		{Name: "Name", Value: install.Name},
		{Name: "Version", Value: install.Version},
		{Name: "Dry run", Value: YesNo(install.DryRun)},
	}
	for _, fetched := range install.Fetched {
		record = append(record, Field{Name: "Fetched", Value: fetched})
	}
	for _, kept := range install.Kept {
		record = append(record, Field{Name: "Kept", Value: kept})
	}
	keys := make([]string, 0, len(install.Components))
	for key := range install.Components {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		record = append(record, Field{Name: "Components " + key, Value: JSONText(install.Components[key])})
	}
	return record
}
