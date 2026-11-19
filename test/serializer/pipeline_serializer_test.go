package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestPipelineSerializer_ModelToTable_MapsFields(t *testing.T) {
	software_id := 9
	status := "success"
	table := cli.PipelineSerializer{}.ModelToTable([]domain.PipelineModel{{ID: 1, FullName: "games/pong", Platform: "love", Active: true, SoftwareID: &software_id, LastPipelineStatus: &status}})

	assert.Equal(t, []string{"1", "games/pong", "love", "yes", "9", "success", "-"}, table.Rows[0])
}

func TestPipelineSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.PipelineSerializer{}.ModelToTable([]domain.PipelineModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestPipelineSerializer_RequestToAttributes_SoftwareEmptyIsNull(t *testing.T) {
	serializer := cli.PipelineSerializer{}
	flags := parse(t, serializer.Fields(), "--software", "")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": nil}, attributes)
}

func TestPipelineSerializer_SyncToTable_OneRowPerChange(t *testing.T) {
	table := cli.PipelineSerializer{}.SyncToTable(domain.PipelineSyncModel{
		Created: []domain.PipelineModel{{ID: 1}}, Updated: []domain.PipelineModel{{ID: 2}}, Deactivated: []domain.PipelineModel{{ID: 3}},
	})

	assert.Len(t, table.Rows, 3)
	assert.Equal(t, []string{"1", "created"}, table.Rows[0][:2])
	assert.Equal(t, []string{"3", "deactivated"}, table.Rows[2][:2])
}

func TestPipelineSerializer_RunsToTable_ShortCommitFirstLine(t *testing.T) {
	commit := "141d95641234567890"
	message := "Release 1.0\nmore"
	table := cli.PipelineSerializer{}.RunsToTable([]domain.PipelineRunModel{{Number: 5, Commit: &commit, Message: &message}})

	assert.Equal(t, []string{"5", "-", "-", "-", "141d9564", "-", "Release 1.0"}, table.Rows[0])
}

func TestPipelineSerializer_RunsToTable_Empty(t *testing.T) {
	assert.Empty(t, cli.PipelineSerializer{}.RunsToTable(nil).Rows)
}
