package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestPlatformLinkSerializer_RequestToAttributes_RejectsUnknownPlatform(t *testing.T) {
	serializer := cli.PlatformLinkSerializer{}
	flags := parse(t, serializer.Fields(), "--platform", "atari")

	_, err := serializer.RequestToAttributes(flags)

	assert.ErrorContains(t, err, "--platform")
}

func TestPlatformLinkSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.PlatformLinkSerializer{}.ModelToTable([]domain.PlatformLinkModel{
		{ID: 3, Platform: "c64", Name: "VICE", URL: "https://vice.example.org", Position: 1},
	})

	assert.Equal(t, []string{"3", "c64", "1", "VICE", "https://vice.example.org"}, table.Rows[0])
}

func TestPlatformLinkSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.PlatformLinkSerializer{}.ModelToTable([]domain.PlatformLinkModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestPlatformLinkSerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.PlatformLinkSerializer{}.ModelToRecord(domain.PlatformLinkModel{ID: 3})

	assert.Equal(t, "ID", record[0].Name)
	assert.Equal(t, "3", record[0].Value)
}
