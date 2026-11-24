package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// MakefileSerializer renders a [domain.MakefileReportModel].
type MakefileSerializer struct{}

// Heading names the project and its platform: "./pong: platform godot".
func (MakefileSerializer) Heading(report domain.MakefileReportModel) string {
	text := report.Directory + ": platform " + report.Platform
	if report.App != "" {
		text += ", app " + report.App
	}
	return text
}

// ReportToTable renders the targets, the name first.
func (MakefileSerializer) ReportToTable(report domain.MakefileReportModel) Table {
	table := Table{Headers: []string{"Target", "Status", "Required because", "Detail"}}
	for _, target := range report.Targets {
		table.Rows = append(table.Rows, []string{target.Name, target.Status, target.Reason, target.Detail})
	}
	return table
}
