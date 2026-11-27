package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareForkSerializer_RequestToAttributes_RejectsUnknownStatus(t *testing.T) {
	serializer := cli.SoftwareForkSerializer{}
	flags := parse(t, serializer.Fields(), "--status", "approved")

	_, err := serializer.RequestToAttributes(flags)

	assert.ErrorContains(t, err, "--status")
}

func TestSoftwareForkSerializer_ModelToTable_MapsFields(t *testing.T) {
	name, title := "pong", "Pong Plus"
	table := cli.SoftwareForkSerializer{}.ModelToTable([]domain.SoftwareForkModel{
		{ID: 4, SoftwareName: &name, EngineURL: "https://other.example.org", Name: "pong-plus", Title: &title, Status: "pending"},
	})

	assert.Equal(t, []string{"4", "pending", "pong", "https://other.example.org", "pong-plus", "Pong Plus", "-"}, table.Rows[0])
}

func TestSoftwareForkSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.SoftwareForkSerializer{}.ModelToTable([]domain.SoftwareForkModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestSoftwareForkSerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.SoftwareForkSerializer{}.ModelToRecord(domain.SoftwareForkModel{ID: 4})

	assert.Equal(t, "ID", record[0].Name)
	assert.Equal(t, "4", record[0].Value)
}

func TestSoftwareSourceSerializer_ModelToRecord_ShowsRefAndPublic(t *testing.T) {
	ref, public := "main", true
	record := cli.SoftwareSourceSerializer{}.ModelToRecord(domain.SoftwareSourceModel{ID: 6, URL: "https://git.example.org/pong", Ref: &ref, Public: &public})

	assert.Equal(t, cli.Field{Name: "ID", Value: "6"}, record[0])
	assert.Contains(t, record, cli.Field{Name: "Ref", Value: "main"})
	assert.Contains(t, record, cli.Field{Name: "Public", Value: "yes"})
}

func TestSoftwareLinkSerializer_Fields_HaveNoKind(t *testing.T) {
	serializer := cli.SoftwareLinkSerializer{}
	flags := parse(t, serializer.Fields(), "--label", "Store", "--url", "https://store.example.org/pong")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"label": "Store", "url": "https://store.example.org/pong"}, attributes)
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
