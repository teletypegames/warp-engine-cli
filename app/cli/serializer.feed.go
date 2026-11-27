package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// FeedSerializer renders the feeds and their checks as tables and
// records.
type FeedSerializer struct{}

// ModelToTable renders feeds as a table: key, type, platforms, whether
// it is materialized and served, description. The key is the first column,
// the one -o id prints.
func (s FeedSerializer) ModelToTable(feeds []domain.FeedModel) Table {
	table := Table{Headers: []string{"Key", "Type", "Platforms", "Static", "Dynamic", "Description"}}
	for _, feed := range feeds {
		table.Rows = append(table.Rows, []string{
			feed.Key, feed.Type, List(feed.Platforms), YesNo(feed.Static), YesNo(feed.Dynamic),
			Truncate(feed.Description, 60),
		})
	}
	return table
}

// CheckToRecord renders the check of one feed: a line per title on offer
// ("ok") and per software left out ("skipped", with the reason).
func (s FeedSerializer) CheckToRecord(check domain.FeedCheckModel) Record {
	record := Record{{Name: "Feed", Value: check.Key + " (" + check.Type + ")"}}
	if !check.Selection {
		return append(record, Field{Name: "Selection", Value: "decided per caller"})
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
