package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestRuntimeSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.RuntimeSerializer{}.ModelToTable([]domain.RuntimeVersionModel{{ID: 1, Name: "emulatorjs", Version: "4.2.3", Source: "download", SizeBytes: 512, Releases: 2}})

	assert.Equal(t, []string{"1", "emulatorjs", "4.2.3", "download", "512 B", "2", "-"}, table.Rows[0])
}

func TestRuntimeSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.RuntimeSerializer{}.ModelToTable([]domain.RuntimeVersionModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestRuntimeSerializer_ModelToRecord_ComponentLines(t *testing.T) {
	record := cli.RuntimeSerializer{}.ModelToRecord(domain.RuntimeVersionModel{ID: 1, Components: []domain.RuntimeComponentModel{{Kind: "core", Key: "vice_x64sc"}}})

	assert.Equal(t, "1", record[0].Value)
	assert.Equal(t, cli.Field{Name: "Component", Value: "core vice_x64sc"}, record[len(record)-1])
}

func TestRuntimeSerializer_InstallToRecord_ListsFiles(t *testing.T) {
	record := cli.RuntimeSerializer{}.InstallToRecord(domain.RuntimeInstallModel{
		Name: "tic80", Version: "1.2", DryRun: true, Fetched: []string{"a.js"}, Kept: []string{"b.js"},
		Components: map[string]any{"core": []any{"x"}},
	})

	assert.Contains(t, record, cli.Field{Name: "Dry run", Value: "yes"})
	assert.Contains(t, record, cli.Field{Name: "Fetched", Value: "a.js"})
	assert.Contains(t, record, cli.Field{Name: "Kept", Value: "b.js"})
	assert.Contains(t, record, cli.Field{Name: "Components core", Value: `["x"]`})
}
