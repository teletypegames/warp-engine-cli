package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Output formats of the --output flag.
const (
	// OutputTable prints aligned columns for a list and "Field  value" lines
	// for one record; the default.
	OutputTable = "table"
	// OutputJSON prints the records as the API sent them, indented: a JSON
	// array for a list, an object for one record.
	OutputJSON = "json"
	// OutputID prints only the identifier of each record, one per line, for
	// shell loops.
	OutputID = "id"
)

// Outputs lists the values --output takes.
var Outputs = []string{OutputTable, OutputJSON, OutputID}

// Table is a list rendered for [OutputTable]: a header row and one row per
// record. The first column is the identifier, which [OutputID] prints alone.
type Table struct {
	// Headers are the column titles.
	Headers []string
	// Rows are the cells, one slice per record, as long as Headers.
	Rows [][]string
}

// Field is one line of a [Record]: a label and its value.
type Field struct {
	// Name is the label, "Title".
	Name string
	// Value is the rendered value; [Text] and friends render the API's types.
	Value string
}

// Record is one record rendered for [OutputTable]: its fields in order. The
// first field is the identifier, which [OutputID] prints alone.
type Record []Field

// outputOf returns the --output flag of cmd, [OutputTable] when it has none.
func outputOf(cmd *cobra.Command) string {
	if flag := cmd.Flags().Lookup("output"); flag != nil && flag.Value.String() != "" {
		return flag.Value.String()
	}
	return OutputTable
}

// PrintList prints a list in the format of --output: value (the records) as
// JSON, or table as columns, or the first column of table.
func PrintList(cmd *cobra.Command, value any, table Table) error {
	writer := cmd.OutOrStdout()
	switch outputOf(cmd) {
	case OutputJSON:
		return writeJSON(writer, value)
	case OutputID:
		for _, row := range table.Rows {
			if len(row) > 0 {
				fmt.Fprintln(writer, row[0])
			}
		}
		return nil
	default:
		return writeTable(writer, table)
	}
}

// PrintRecord prints one record in the format of --output: value as JSON,
// record as aligned "Name  value" lines, or the first field's value.
func PrintRecord(cmd *cobra.Command, value any, record Record) error {
	writer := cmd.OutOrStdout()
	switch outputOf(cmd) {
	case OutputJSON:
		return writeJSON(writer, value)
	case OutputID:
		if len(record) > 0 {
			fmt.Fprintln(writer, record[0].Value)
		}
		return nil
	default:
		tab := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
		for _, field := range record {
			value := strings.ReplaceAll(field.Value, "\n", "\n\t")
			fmt.Fprintf(tab, "%s\t%s\n", field.Name, value)
		}
		return tab.Flush()
	}
}

// PrintPaging tells, on standard error and only for [OutputTable], that a
// list has more pages than the one printed: "Page 1 of 3 (60 in all): --page 2
// or --all for more". It prints nothing for the last page.
func PrintPaging(cmd *cobra.Command, page int, total_pages int, total int) {
	if outputOf(cmd) != OutputTable || page < 1 || page >= total_pages {
		return
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "Page %d of %d (%d in all): --page %d or --all for more\n", page, total_pages, total, page+1)
}

// PrintDone prints a short confirmation (for a delete, a trigger …) on
// standard output for [OutputTable] and nothing otherwise, so the JSON and
// ID outputs stay clean.
func PrintDone(cmd *cobra.Command, format string, args ...any) {
	if outputOf(cmd) == OutputTable {
		fmt.Fprintf(cmd.OutOrStdout(), format+"\n", args...)
	}
}

// writeJSON writes value as indented JSON.
func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

// writeTable writes table as aligned columns, or "(none)" when it has no rows.
func writeTable(writer io.Writer, table Table) error {
	if len(table.Rows) == 0 {
		_, err := fmt.Fprintln(writer, "(none)")
		return err
	}
	tab := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tab, strings.Join(upper(table.Headers), "\t"))
	for _, row := range table.Rows {
		cells := make([]string, len(row))
		for i, cell := range row {
			cells[i] = strings.ReplaceAll(cell, "\n", " ")
		}
		fmt.Fprintln(tab, strings.Join(cells, "\t"))
	}
	return tab.Flush()
}

// upper returns headers in capitals, the way the table prints them.
func upper(headers []string) []string {
	result := make([]string, len(headers))
	for i, header := range headers {
		result[i] = strings.ToUpper(header)
	}
	return result
}

// Text renders an optional string: its value, or "-" for nil or empty.
func Text(value *string) string {
	if value == nil || *value == "" {
		return "-"
	}
	return *value
}

// Number renders an optional integer: its value, or "-" for nil.
func Number(value *int) string {
	if value == nil {
		return "-"
	}
	return strconv.Itoa(*value)
}

// Flag renders an optional boolean: "yes", "no", or "-" for nil.
func Flag(value *bool) string {
	if value == nil {
		return "-"
	}
	return YesNo(*value)
}

// YesNo renders a boolean as "yes" or "no".
func YesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// List renders a list of words comma-separated, or "-" when it is empty.
func List(values []string) string {
	if len(values) == 0 {
		return "-"
	}
	return strings.Join(values, ", ")
}

// Truncate shortens text to at most width runes, ending in "…" when cut.
func Truncate(text string, width int) string {
	runes := []rune(text)
	if len(runes) <= width || width < 2 {
		return text
	}
	return string(runes[:width-1]) + "…"
}

// Bytes renders a size in bytes with a binary unit: "512 B", "1.5 KiB", "3.2 MiB".
func Bytes(size int64) string {
	if size < 1024 {
		return strconv.FormatInt(size, 10) + " B"
	}
	value := float64(size)
	for _, unit := range []string{"KiB", "MiB", "GiB", "TiB"} {
		value /= 1024
		if value < 1024 {
			return strconv.FormatFloat(value, 'f', 1, 64) + " " + unit
		}
	}
	return strconv.FormatFloat(value, 'f', 1, 64) + " PiB"
}

// JSONText renders any value as compact JSON, for the free-form fields
// (transport options, results) in a table; nil is "-".
func JSONText(value any) string {
	if value == nil {
		return "-"
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(data)
}
