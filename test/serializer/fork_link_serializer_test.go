package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestForkLinkSerializer_RequestToAttributes_RejectsUnknownStatus(t *testing.T) {
	serializer := cli.ForkLinkSerializer{}
	flags := parse(t, serializer.Fields(), "--status", "approved")

	_, err := serializer.RequestToAttributes(flags)

	assert.ErrorContains(t, err, "--status")
}

func TestForkLinkSerializer_ModelToTable_MapsFields(t *testing.T) {
	name, title := "pong", "Pong Plus"
	table := cli.ForkLinkSerializer{}.ModelToTable([]domain.ForkLinkModel{
		{ID: 4, SoftwareName: &name, ChildEngine: "https://other.example.org", ChildName: "pong-plus", ChildTitle: &title, Status: "pending"},
	})

	assert.Equal(t, []string{"4", "pending", "pong", "https://other.example.org", "pong-plus", "Pong Plus", "-"}, table.Rows[0])
}

func TestForkLinkSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.ForkLinkSerializer{}.ModelToTable([]domain.ForkLinkModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestForkLinkSerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.ForkLinkSerializer{}.ModelToRecord(domain.ForkLinkModel{ID: 4})

	assert.Equal(t, "ID", record[0].Name)
	assert.Equal(t, "4", record[0].Value)
}

func TestExternalLinkSerializer_ModelToTable_ShowsKindRefAndPublic(t *testing.T) {
	ref, public := "main", true
	table := cli.ExternalLinkSerializer{}.ModelToTable([]domain.ExternalLinkModel{
		{ID: 6, Kind: "source", Label: "Source", URL: "https://git.example.org/pong", Ref: &ref, Public: &public},
	})

	assert.Equal(t, []string{"6", "source", "Source", "https://git.example.org/pong", "main", "yes"}, table.Rows[0])
}

func TestExternalLinkSerializer_RequestToAttributes_RejectsUnknownKind(t *testing.T) {
	serializer := cli.ExternalLinkSerializer{}
	flags := parse(t, serializer.Fields(), "--kind", "video")

	_, err := serializer.RequestToAttributes(flags)

	assert.ErrorContains(t, err, "--kind")
}

func TestSoftwareSerializer_ScaffoldToAttributes_MergesForkOfOverData(t *testing.T) {
	flags := parse(t, cli.SoftwareSerializer{}.ScaffoldFields(), "--data", `{"forkOf":{"engine":"https://old.example.org","version":"1.0.0"}}`,
		"--fork-of-engine", "https://games.example.org", "--fork-of-name", "pong")

	attributes, err := cli.SoftwareSerializer{}.ScaffoldToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"forkOf": map[string]any{"engine": "https://games.example.org", "name": "pong", "version": "1.0.0"}}, attributes)
}

func TestSoftwareSerializer_ScaffoldToAttributes_NoForkOfSendsNone(t *testing.T) {
	flags := parse(t, cli.SoftwareSerializer{}.ScaffoldFields(), "--name", "my-game")

	attributes, err := cli.SoftwareSerializer{}.ScaffoldToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"name": "my-game"}, attributes)
}
