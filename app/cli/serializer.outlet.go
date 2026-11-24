package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// OutletSerializer renders the outlets and their checks as tables and
// records.
type OutletSerializer struct{}

// ModelToTable renders outlets as a table: key, type, platforms, whether
// it is materialized and served, description. The key is the first column,
// the one -o id prints.
func (s OutletSerializer) ModelToTable(outlets []domain.OutletModel) Table {
	table := Table{Headers: []string{"Key", "Type", "Platforms", "Materialized", "Served", "Description"}}
	for _, outlet := range outlets {
		table.Rows = append(table.Rows, []string{
			outlet.Key, outlet.Type, List(outlet.Platforms), YesNo(outlet.Materialized), YesNo(outlet.Served),
			Truncate(outlet.Description, 60),
		})
	}
	return table
}

// CheckToRecord renders the check of one outlet: a line per title on offer
// ("ok") and per software left out ("skipped", with the reason).
func (s OutletSerializer) CheckToRecord(check domain.OutletCheckModel) Record {
	record := Record{{Name: "Outlet", Value: check.Key + " (" + check.Type + ")"}}
	if !check.Offer {
		return append(record, Field{Name: "Offer", Value: "decided per caller"})
	}
	record = append(record, Field{Name: "On offer", Value: strconv.Itoa(len(check.Items))},
		Field{Name: "Skipped", Value: strconv.Itoa(len(check.Skipped))})
	for _, item := range check.Items {
		record = append(record, Field{Name: "  ok " + item.Name, Value: item.Version + " (" + item.Kind + ")"})
	}
	for _, skip := range check.Skipped {
		record = append(record, Field{Name: "  skipped " + skip.Name, Value: skip.Reason})
	}
	return record
}
