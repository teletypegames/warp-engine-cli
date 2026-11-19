package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestLocalComponentSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.LocalComponentSerializer{}.ModelToTable([]domain.LocalComponentModel{{Key: "woodpecker", Label: "CI", Role: "ci", State: "ready", Requires: []string{"forge"}}})

	assert.Equal(t, []string{"woodpecker", "CI", "ci", "ready", "forge", "-"}, table.Rows[0])
}

func TestLocalComponentSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.LocalComponentSerializer{}.ModelToTable([]domain.LocalComponentModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestLocalOperationSerializer_ModelToRecord_StepsAndResult(t *testing.T) {
	record := cli.LocalOperationSerializer{}.ModelToRecord(domain.LocalOperationModel{
		ID: 3, Kind: "local_ci", State: "succeeded",
		Steps:  []domain.LocalOperationStepModel{{At: "soon", Text: "create the repository"}},
		Result: map[string]any{"pushUrl": "http://x/pong.git", "pipelineId": 4.0},
	})

	assert.Equal(t, cli.Field{Name: "ID", Value: "3"}, record[0])
	assert.Contains(t, record, cli.Field{Name: "Step", Value: "soon  create the repository"})
	assert.Contains(t, record, cli.Field{Name: "Result pushUrl", Value: "http://x/pong.git"})
	assert.Contains(t, record, cli.Field{Name: "Result pipelineId", Value: "4"})
}
