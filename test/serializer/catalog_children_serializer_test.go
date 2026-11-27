package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestReleaseSerializer_ModelToTable_MapsFields(t *testing.T) {
	sha := "3fb3ed8aa1b2c3d4e5f6"
	created := "2026-11-19T10:00:00Z"
	table := cli.ReleaseSerializer{}.ModelToTable([]domain.ReleaseModel{{
		ID: 42, Version: "1.2.0", Channel: "stable", CommitSha: &sha, CreatedAt: &created,
		Assets: []domain.ReleaseAssetModel{{ID: 1}, {ID: 2}},
	}})

	assert.Equal(t, []string{"42", "1.2.0", "stable", "3fb3ed8aa1", "2", created}, table.Rows[0])
}

func TestReleaseSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.ReleaseSerializer{}.ModelToTable(nil)

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestReleaseSerializer_ModelToRecord_IDFirstAndAssets(t *testing.T) {
	record := cli.ReleaseSerializer{}.ModelToRecord(domain.ReleaseModel{ID: 42, Assets: []domain.ReleaseAssetModel{{ID: 7, Kind: "html", Path: "/p"}}})

	assert.Equal(t, cli.Field{Name: "ID", Value: "42"}, record[0])
	assert.Equal(t, cli.Field{Name: "Asset", Value: "#7 html /p"}, record[len(record)-1])
}

func TestReleaseSerializer_FlagsToFilter_OnlyGivenFilters(t *testing.T) {
	serializer := cli.ReleaseSerializer{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFilters(flags, serializer.Filters())
	assert.NoError(t, flags.Parse([]string{"--commit-sha", "abc"}))

	assert.Equal(t, domain.Filter{"commitSha": "abc"}, serializer.FlagsToFilter(flags))
}

func TestReleaseAssetSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.ReleaseAssetSerializer{}.ModelToTable([]domain.ReleaseAssetModel{{ID: 7, ReleaseID: 42, Kind: "html", Path: "/p/index.html"}})

	assert.Equal(t, []string{"7", "42", "html", "/p/index.html"}, table.Rows[0])
}

func TestReleaseAssetSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.ReleaseAssetSerializer{}.ModelToTable([]domain.ReleaseAssetModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestReleaseAssetSerializer_RequestToAttributes_PassesAnyKind(t *testing.T) {
	serializer := cli.ReleaseAssetSerializer{}
	flags := parse(t, serializer.Fields(), "--kind", "cli_linux_x64")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"kind": "cli_linux_x64"}, attributes)
}

func TestSoftwareLinkSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.SoftwareLinkSerializer{}.ModelToTable([]domain.SoftwareLinkModel{{ID: 2, Label: "Store", URL: "https://x"}})

	assert.Equal(t, []string{"2", "Store", "https://x"}, table.Rows[0])
}

func TestSoftwareLinkSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.SoftwareLinkSerializer{}.ModelToTable(nil)

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestSoftwareLinkSerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.SoftwareLinkSerializer{}.ModelToRecord(domain.SoftwareLinkModel{ID: 2})

	assert.Equal(t, cli.Field{Name: "ID", Value: "2"}, record[0])
}

func TestSoftwareImageSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.SoftwareImageSerializer{}.ModelToTable([]domain.SoftwareImageModel{{ID: 31, ImageID: 118, Position: 1, Origin: "sync"}})

	assert.Equal(t, []string{"31", "118", "1", "no", "sync", "-"}, table.Rows[0])
}

func TestSoftwareImageSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.SoftwareImageSerializer{}.ModelToTable(nil)

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestSoftwareImageSerializer_CreateToAttributes_TypedValues(t *testing.T) {
	serializer := cli.SoftwareImageSerializer{}
	flags := parse(t, serializer.CreateFields(), "--image-id", "118", "--default=false", "--position", "0")

	attributes, err := serializer.CreateToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"imageId": 118, "isDefault": false, "position": 0}, attributes)
}

func TestExtraSerializer_ModelToRecord_AssociationFirstValueAsJSON(t *testing.T) {
	record := cli.ExtraSerializer{}.ModelToRecord(domain.ExtraModel{Association: "seo_metadata", Value: map[string]any{"title": "Pong"}})

	assert.Equal(t, cli.Field{Name: "Association", Value: "seo_metadata"}, record[0])
	assert.Equal(t, "{\n  \"title\": \"Pong\"\n}", record[2].Value)
}

func TestExtraSerializer_RequestToAttributes_DataWithoutValue(t *testing.T) {
	serializer := cli.ExtraSerializer{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	serializer.DeclareFlags(flags)
	assert.NoError(t, flags.Parse([]string{"--data", `{"value":{"title":"Pong"}}`}))

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"value": map[string]any{"title": "Pong"}}, attributes)
}
