package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestCiRepositorySerializer_ModelToTable_MapsFields(t *testing.T) {
	software_id := 9
	status := "success"
	table := cli.CiRepositorySerializer{}.ModelToTable([]domain.CiRepositoryModel{{ID: 1, FullName: "games/pong", Platform: "love", Active: true, SoftwareID: &software_id, LastPipelineStatus: &status}})

	assert.Equal(t, []string{"1", "games/pong", "love", "yes", "9", "success", "-"}, table.Rows[0])
}

func TestCiRepositorySerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.CiRepositorySerializer{}.ModelToTable([]domain.CiRepositoryModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestCiRepositorySerializer_RequestToAttributes_SoftwareEmptyIsNull(t *testing.T) {
	serializer := cli.CiRepositorySerializer{}
	flags := parse(t, serializer.Fields(), "--software", "")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"softwareId": nil}, attributes)
}

func TestCiRepositorySerializer_SyncToTable_OneRowPerChange(t *testing.T) {
	table := cli.CiRepositorySerializer{}.SyncToTable(domain.CiRepositorySyncModel{
		Created: []domain.CiRepositoryModel{{ID: 1}}, Updated: []domain.CiRepositoryModel{{ID: 2}}, Deactivated: []domain.CiRepositoryModel{{ID: 3}},
	})

	assert.Len(t, table.Rows, 3)
	assert.Equal(t, []string{"1", "created"}, table.Rows[0][:2])
	assert.Equal(t, []string{"3", "deactivated"}, table.Rows[2][:2])
}

func TestCiRepositorySerializer_RunsToTable_ShortCommitFirstLine(t *testing.T) {
	commit := "141d95641234567890"
	message := "Release 1.0\nmore"
	table := cli.CiRepositorySerializer{}.RunsToTable([]domain.PipelineRunModel{{Number: 5, Commit: &commit, Message: &message}})

	assert.Equal(t, []string{"5", "-", "-", "-", "141d9564", "-", "Release 1.0"}, table.Rows[0])
}

func TestCiRepositorySerializer_RunsToTable_Empty(t *testing.T) {
	assert.Empty(t, cli.CiRepositorySerializer{}.RunsToTable(nil).Rows)
}
