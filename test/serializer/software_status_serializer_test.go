package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestSoftwareStatusSerializer_RequestToAttributes_TypesTheFields(t *testing.T) {
	serializer := cli.SoftwareStatusSerializer{}
	flags := parse(t, serializer.Fields(), "--position", "2", "--default-visible=false", "--badge", "")

	attributes, err := serializer.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"position": 2, "defaultVisible": false, "badge": nil}, attributes)
}

func TestSoftwareStatusSerializer_ModelToTable_MapsFields(t *testing.T) {
	badge := "NEW"
	table := cli.SoftwareStatusSerializer{}.ModelToTable([]domain.SoftwareStatusModel{
		{ID: 2, Key: "released", Label: "Released", Position: 1, DefaultVisible: true, Badge: &badge, SoftwareCount: 14},
	})

	assert.Equal(t, []string{"2", "released", "Released", "1", "yes", "NEW", "14"}, table.Rows[0])
}

func TestSoftwareStatusSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.SoftwareStatusSerializer{}.ModelToTable([]domain.SoftwareStatusModel{})

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestSoftwareStatusSerializer_ModelToRecord_IDFirst(t *testing.T) {
	record := cli.SoftwareStatusSerializer{}.ModelToRecord(domain.SoftwareStatusModel{ID: 2, Key: "released"})

	assert.Equal(t, "ID", record[0].Name)
	assert.Equal(t, "2", record[0].Value)
}
