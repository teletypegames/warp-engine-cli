package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// ReleaseSerializer turns the flags of the release commands into the API's
// terms and the releases into tables and records.
type ReleaseSerializer struct{}

// ReleaseSorts are the orders "release list --sort" takes; the server's
// default is -createdAt.
var ReleaseSorts = []string{"createdAt", "version"}

// Filters are the filter flags of "release list".
func (s ReleaseSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "commit-sha", Key: "commitSha", Usage: "built from this commit"},
		{Flag: "channel", Key: "channel", Usage: "channel is: stable, beta, dev"},
	}
}

// Fields are the field flags of "release create" and "release update". The
// API takes only the version: the channel follows from it (dev-… is dev,
// 1.2.0-beta.1 is beta, anything else stable) and the commit is set by the CI.
func (s ReleaseSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "version", Key: "version", Usage: "version, \"1.2.0\"; it sets the channel"},
	}
}

// FlagsToFilter reads the filter flags of "release list".
func (s ReleaseSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of create and update.
func (s ReleaseSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders releases as a table: ID, version, channel, commit
// (shortened), number of files, created.
func (s ReleaseSerializer) ModelToTable(releases []domain.ReleaseModel) Table {
	table := Table{Headers: []string{"ID", "Version", "Channel", "Commit", "Files", "Created"}}
	for _, release := range releases {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(release.ID), release.Version, release.Channel, shortSha(release.CommitSha),
			strconv.Itoa(len(release.Assets)), Text(release.CreatedAt),
		})
	}
	return table
}

// ModelToRecord renders one release with every field and a line per asset.
func (s ReleaseSerializer) ModelToRecord(release domain.ReleaseModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(release.ID)},
		{Name: "Software", Value: strconv.Itoa(release.SoftwareID)},
		{Name: "Version", Value: release.Version},
		{Name: "Channel", Value: release.Channel},
		{Name: "Commit", Value: Text(release.CommitSha)},
		{Name: "Created", Value: Text(release.CreatedAt)},
		{Name: "Updated", Value: Text(release.UpdatedAt)},
	}
	for _, asset := range release.Assets {
		record = append(record, Field{Name: "Asset", Value: "#" + strconv.Itoa(asset.ID) + " " + asset.Kind + " " + asset.Path})
	}
	return record
}

// shortSha renders a commit as its first 10 characters, or "-" for none.
func shortSha(sha *string) string {
	if sha == nil || len(*sha) <= 10 {
		return Text(sha)
	}
	return (*sha)[:10]
}
