package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestLocalProviderSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.LocalProviderSerializer{}.ModelToTable([]domain.LocalProviderModel{{Key: "woodpecker", Label: "CI", Role: "ci", State: "ready", Requires: []string{"forge"}}})

	assert.Equal(t, []string{"woodpecker", "CI", "ci", "ready", "forge", "-"}, table.Rows[0])
}

func TestLocalProviderSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.LocalProviderSerializer{}.ModelToTable([]domain.LocalProviderModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestLocalTaskSerializer_ModelToRecord_StepsAndResult(t *testing.T) {
	record := cli.LocalTaskSerializer{}.ModelToRecord(domain.LocalTaskModel{
		ID: 3, Kind: "local_ci", State: "succeeded",
		Steps:  []domain.LocalTaskStepModel{{At: "soon", Text: "create the repository"}},
		Result: map[string]any{"pushUrl": "http://x/pong.git", "ciRepositoryId": 4.0},
	})

	assert.Equal(t, cli.Field{Name: "ID", Value: "3"}, record[0])
	assert.Contains(t, record, cli.Field{Name: "Step", Value: "soon  create the repository"})
	assert.Contains(t, record, cli.Field{Name: "Result pushUrl", Value: "http://x/pong.git"})
	assert.Contains(t, record, cli.Field{Name: "Result ciRepositoryId", Value: "4"})
}
