package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestFileSerializer_ModelToTable_FileAndDirectory(t *testing.T) {
	size := int64(1536)
	downloads := 4
	table := cli.FileSerializer{}.ModelToTable([]domain.FileEntryModel{
		{Name: "pong", Path: "pong", Type: "directory"},
		{Name: "a.zip", Path: "pong/a.zip", Type: "file", Size: &size, Downloads: &downloads},
	})

	assert.Equal(t, []string{"pong/", "directory", "-", "-", "-"}, table.Rows[0])
	assert.Equal(t, []string{"pong/a.zip", "file", "1.5 KiB", "-", "4"}, table.Rows[1])
}

func TestFileSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.FileSerializer{}.ModelToTable([]domain.FileEntryModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}
