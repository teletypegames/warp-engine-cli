package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareCategorySerializer_RequestToAttributes_OnlyChangedFlags(t *testing.T) {
	serializer := cli.SoftwareCategorySerializer{}
	flags := parse(t, serializer.Fields(), "--key", "tool", "--needs-review")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"key": "tool", "needsReview": true}, attributes)
}

func TestSoftwareCategorySerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.SoftwareCategorySerializer{}.ModelToTable([]domain.SoftwareCategoryModel{
		{ID: 5, Key: "puzzle", Label: "Puzzle", Position: 3, NeedsReview: true, SoftwareCount: 2},
	})

	assert.Equal(t, []string{"5", "puzzle", "Puzzle", "3", "yes", "2"}, table.Rows[0])
}

func TestSoftwareCategorySerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.SoftwareCategorySerializer{}.ModelToTable(nil)

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestSoftwareCategorySerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.SoftwareCategorySerializer{}.ModelToRecord(domain.SoftwareCategoryModel{ID: 5})

	assert.Equal(t, "ID", record[0].Name)
	assert.Equal(t, "5", record[0].Value)
}
