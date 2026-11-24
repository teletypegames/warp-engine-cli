package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// BuildSerializer turns the flags of "release upload" and "release
// publish" into a [domain.BuildTargetModel] and the answers into output.
type BuildSerializer struct{}

// DeclareTarget adds --name and --version, and with publish --platform,
// --app and --commit.
func (BuildSerializer) DeclareTarget(flags *pflag.FlagSet, publish bool) {
	flags.String("name", "", "the software's name in the catalog (required)")
	flags.String("version", "", "the release's version (required)")
	if publish {
		flags.String("platform", "", "the software's platform (required)")
		flags.String("app", "", "the platform variant (golang: cli or gui); the default one without it")
		flags.String("commit", "", "the commit the release was built from")
	}
}

// FlagsToTarget reads the flags [BuildSerializer.DeclareTarget] added.
func (BuildSerializer) FlagsToTarget(flags *pflag.FlagSet) domain.BuildTargetModel {
	value := func(name string) string {
		text, _ := flags.GetString(name)
		return text
	}
	return domain.BuildTargetModel{
		Name: value("name"), Version: value("version"), Platform: value("platform"), App: value("app"), Commit: value("commit"),
	}
}

// UploadsToTable renders uploaded files: the stored name first.
func (BuildSerializer) UploadsToTable(uploads []domain.BuildUploadModel) Table {
	table := Table{Headers: []string{"File", "Size", "SHA256"}}
	for _, upload := range uploads {
		table.Rows = append(table.Rows, []string{upload.File, Bytes(upload.Size), upload.SHA256})
	}
	return table
}

// PublishToRecord renders the answer of a publish, the name first.
func (BuildSerializer) PublishToRecord(published domain.BuildPublishModel) Record {
	return Record{
		{Name: "Name", Value: published.Name},
		{Name: "Platform", Value: published.Platform},
		{Name: "Version", Value: published.Version},
		{Name: "Published", Value: strconv.FormatBool(published.Published)},
	}
}
