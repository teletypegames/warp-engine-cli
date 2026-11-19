package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// parse declares fields on a new flag set and parses args into it.
func parse(t *testing.T, fields []cli.FieldFlag, args ...string) *pflag.FlagSet {
	t.Helper()
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFields(flags, fields)
	assert.NoError(t, flags.Parse(args))
	return flags
}

func TestSoftwareSerializer_RequestToAttributes_OnlyChangedFlags(t *testing.T) {
	serializer := cli.SoftwareSerializer{}
	flags := parse(t, serializer.Fields(), "--title", "Pong", "--shelf-name", "PONG")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"title": "Pong", "shelfName": "PONG"}, attributes)
}

func TestSoftwareSerializer_RequestToAttributes_DataIsMergedUnderFlags(t *testing.T) {
	serializer := cli.SoftwareSerializer{}
	flags := parse(t, serializer.Fields(), "--data", `{"title":"Old","story":"Long"}`, "--title", "New")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"title": "New", "story": "Long"}, attributes)
}

func TestSoftwareSerializer_ModelToTable_MapsFields(t *testing.T) {
	status := "released"
	table := cli.SoftwareSerializer{}.ModelToTable([]domain.SoftwareModel{{ID: 4, Name: "pong", Title: "Pong", Platform: "love", Status: &status}})

	assert.Equal(t, []string{"4", "pong", "Pong", "love", "released", "-", "-"}, table.Rows[0])
}

func TestSoftwareSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.SoftwareSerializer{}.ModelToTable([]domain.SoftwareModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestSoftwareSerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.SoftwareSerializer{}.ModelToRecord(domain.SoftwareModel{ID: 4, Name: "pong"})

	assert.Equal(t, "ID", record[0].Name)
	assert.Equal(t, "4", record[0].Value)
}
