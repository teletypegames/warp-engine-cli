package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SoftwareCategorySerializer turns the flags of the category commands into
// the API's terms and the categories into tables and records.
type SoftwareCategorySerializer struct{}

// SoftwareCategorySorts are the orders "category list --sort" takes.
var SoftwareCategorySorts = []string{"position", "key"}

// Filters are the filter flags of "category list".
func (s SoftwareCategorySerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "needs-review", Key: "needsReview", Usage: "waiting for review or not", Bool: true},
	}
}

// Fields are the field flags of "category create" and "category update".
func (s SoftwareCategorySerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "key", Key: "key", Usage: "the word a software's category holds"},
		{Flag: "label", Key: "label", Usage: "display name"},
		{Flag: "position", Key: "position", Kind: KindInt, Usage: "order, lowest first"},
		{Flag: "needs-review", Key: "needsReview", Kind: KindBool, Usage: "still to be reviewed (--needs-review=false approves it)"},
	}
}

// FlagsToFilter reads the filter flags of "category list".
func (s SoftwareCategorySerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of create and update.
func (s SoftwareCategorySerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders categories as a table: ID, key, label, position,
// needs review, number of softwares.
func (s SoftwareCategorySerializer) ModelToTable(categories []domain.SoftwareCategoryModel) Table {
	table := Table{Headers: []string{"ID", "Key", "Label", "Position", "Review", "Softwares"}}
	for _, category := range categories {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(category.ID), category.Key, Truncate(category.Label, 40), strconv.Itoa(category.Position),
			YesNo(category.NeedsReview), strconv.Itoa(category.SoftwareCount),
		})
	}
	return table
}

// ModelToRecord renders one category with every field.
func (s SoftwareCategorySerializer) ModelToRecord(category domain.SoftwareCategoryModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(category.ID)},
		{Name: "Key", Value: category.Key},
		{Name: "Label", Value: category.Label},
		{Name: "Position", Value: strconv.Itoa(category.Position)},
		{Name: "Needs review", Value: YesNo(category.NeedsReview)},
		{Name: "Softwares", Value: strconv.Itoa(category.SoftwareCount)},
		{Name: "Created", Value: Text(category.CreatedAt)},
		{Name: "Updated", Value: Text(category.UpdatedAt)},
	}
}
