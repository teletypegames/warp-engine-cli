package cli

import (
	"sort"
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SoftwareSerializer turns the flags of the software commands into the
// API's terms and the softwares into tables and records.
type SoftwareSerializer struct{}

// SoftwareSorts are the orders "software list --sort" takes.
var SoftwareSorts = []string{"title", "name", "createdAt", "updatedAt"}

// Filters are the filter flags of "software list".
func (s SoftwareSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "name", Key: "name", Usage: "name contains"},
		{Flag: "title", Key: "title", Usage: "title contains"},
		{Flag: "author", Key: "author", Usage: "author contains"},
		{Flag: "platform", Key: "platform", Usage: "platform is"},
		{Flag: "status", Key: "status", Usage: "status key is"},
		{Flag: "category", Key: "category", Usage: "category key is"},
		{Flag: "highlighted", Key: "highlighted", Usage: "highlighted or not", Bool: true},
		{Flag: "query", Key: "q", Usage: "title or name contains"},
	}
}

// Fields are the field flags of "software create" and "software update".
func (s SoftwareSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "name", Key: "name", Usage: "internal name, unique: a-z, 0-9 and dashes"},
		{Flag: "title", Key: "title", Usage: "display name"},
		{Flag: "author", Key: "author", Kind: KindNullable, Usage: "who made it"},
		{Flag: "desc", Key: "desc", Kind: KindNullable, Usage: "short description"},
		{Flag: "story", Key: "story", Kind: KindNullable, Usage: "long description"},
		{Flag: "license", Key: "license", Kind: KindNullable, Usage: "licence"},
		{Flag: "platform", Key: "platform", Usage: "platform", Enum: domain.Platforms},
		{Flag: "status", Key: "status", Kind: KindNullable, Usage: "a status key (see \"status list\")"},
		{Flag: "category", Key: "category", Kind: KindNullable, Usage: "a category key; a new word is registered for review"},
		{Flag: "highlighted", Key: "highlighted", Kind: KindBool, Usage: "show it first in the stores"},
		{Flag: "shelf-name", Key: "shelfName", Kind: KindNullable, Usage: "file name on a shelf (A-Z, 0-9, -)"},
		{Flag: "forks-allowed", Key: "forksAllowed", Kind: KindBool, Usage: "let the developer tools offer a fork (needs a public source link)"},
		{Flag: "fork-of-engine", Key: "forkOfEngine", Kind: KindNullable, Usage: "for a fork: the parent engine's base URL"},
		{Flag: "fork-of-name", Key: "forkOfName", Kind: KindNullable, Usage: "for a fork: the parent's name on its engine"},
		{Flag: "fork-of-version", Key: "forkOfVersion", Kind: KindNullable, Usage: "for a fork: the parent's version"},
		{Flag: "fork-of-repo", Key: "forkOfRepo", Kind: KindNullable, Usage: "for a fork: the parent's source repository"},
	}
}

// ScaffoldFields are the field flags of "software scaffold". The
// --fork-of-… flags have keys under "forkOf." and
// [SoftwareSerializer.ScaffoldToAttributes] gathers them into the forkOf
// object.
func (s SoftwareSerializer) ScaffoldFields() []FieldFlag {
	return []FieldFlag{
		{Flag: "name", Key: "name", Usage: "software and repository name: a-z, 0-9 and dashes"},
		{Flag: "title", Key: "title", Usage: "display name"},
		{Flag: "platform", Key: "platform", Usage: "platform", Enum: domain.Platforms},
		{Flag: "category", Key: "category", Kind: KindNullable, Usage: "a category key"},
		{Flag: "description", Key: "description", Kind: KindNullable, Usage: "short description"},
		{Flag: "license", Key: "license", Kind: KindNullable, Usage: "licence"},
		{Flag: "fork-of-engine", Key: forkOfPrefix + "engine", Usage: "fork a software: the parent engine's base URL"},
		{Flag: "fork-of-name", Key: forkOfPrefix + "name", Usage: "fork a software: the parent's name on its engine"},
		{Flag: "fork-of-version", Key: forkOfPrefix + "version", Usage: "fork a software: the parent's version"},
		{Flag: "fork-of-repo", Key: forkOfPrefix + "repo", Usage: "fork a software: the parent's source repository"},
	}
}

// forkOfPrefix marks the scaffold field keys that belong in the forkOf
// object.
const forkOfPrefix = "forkOf."

// ScaffoldToAttributes reads the field flags of "software scaffold" and
// moves the --fork-of-… values into the forkOf object, merged over a forkOf
// object given in --data.
func (s SoftwareSerializer) ScaffoldToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	attributes, err := FieldsToAttributes(flags, s.ScaffoldFields())
	if err != nil {
		return nil, err
	}
	fork_of := map[string]any{}
	if given, ok := attributes["forkOf"].(map[string]any); ok {
		fork_of = given
	}
	found := false
	for key, value := range attributes {
		if strings.HasPrefix(key, forkOfPrefix) {
			fork_of[strings.TrimPrefix(key, forkOfPrefix)] = value
			delete(attributes, key)
			found = true
		}
	}
	if found {
		attributes["forkOf"] = fork_of
	}
	return attributes, nil
}

// FlagsToFilter reads the filter flags of "software list".
func (s SoftwareSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of create and update.
func (s SoftwareSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// ModelToTable renders softwares as a table: ID, name, title, platform,
// status, category, highlighted.
func (s SoftwareSerializer) ModelToTable(softwares []domain.SoftwareModel) Table {
	table := Table{Headers: []string{"ID", "Name", "Title", "Platform", "Status", "Category", "Highlighted"}}
	for _, software := range softwares {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(software.ID), software.Name, Truncate(software.Title, 40), software.Platform,
			Text(software.Status), Text(software.Category), Flag(software.Highlighted),
		})
	}
	return table
}

// ModelToRecord renders one software with every field, and a count of the
// included releases, images and links when they were asked for.
func (s SoftwareSerializer) ModelToRecord(software domain.SoftwareModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(software.ID)},
		{Name: "Name", Value: software.Name},
		{Name: "Title", Value: software.Title},
		{Name: "Platform", Value: software.Platform},
		{Name: "Author", Value: Text(software.Author)},
		{Name: "Status", Value: Text(software.Status)},
		{Name: "Category", Value: Text(software.Category)},
		{Name: "Highlighted", Value: Flag(software.Highlighted)},
		{Name: "Shelf name", Value: Text(software.ShelfName)},
		{Name: "Forks allowed", Value: YesNo(software.ForksAllowed)},
		{Name: "Fork of engine", Value: Text(software.ForkOfEngine)},
		{Name: "Fork of name", Value: Text(software.ForkOfName)},
		{Name: "Fork of version", Value: Text(software.ForkOfVersion)},
		{Name: "Fork of repo", Value: Text(software.ForkOfRepo)},
		{Name: "License", Value: Text(software.License)},
		{Name: "Description", Value: Text(software.Desc)},
		{Name: "Story", Value: Text(software.Story)},
		{Name: "Created", Value: Text(software.CreatedAt)},
		{Name: "Updated", Value: Text(software.UpdatedAt)},
	}
	for _, release := range software.Releases {
		record = append(record, Field{Name: "Release", Value: "#" + strconv.Itoa(release.ID) + " " + release.Version + " (" + release.Channel + ", " + strconv.Itoa(len(release.Assets)) + " files)"})
	}
	for _, image := range software.Images {
		value := "#" + strconv.Itoa(image.ID) + " " + Text(image.URL)
		if image.IsDefault {
			value += " (default)"
		}
		record = append(record, Field{Name: "Image", Value: value})
	}
	for _, link := range software.ExternalLinks {
		record = append(record, Field{Name: "Link", Value: link.Label + ": " + link.URL})
	}
	keys := make([]string, 0, len(software.Extras))
	for key := range software.Extras {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		record = append(record, Field{Name: "Extra " + key, Value: JSONText(software.Extras[key])})
	}
	return record
}

// ScaffoldToRecord renders the answer of a scaffold: the software, the
// repository's URLs, the pipeline and the warnings.
func (s SoftwareSerializer) ScaffoldToRecord(scaffold domain.ScaffoldModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(scaffold.Software.ID)},
		{Name: "Name", Value: scaffold.Software.Name},
		{Name: "Title", Value: scaffold.Software.Title},
		{Name: "Platform", Value: scaffold.Software.Platform},
		{Name: "Pipeline", Value: Number(scaffold.PipelineID)},
	}
	for _, key := range []string{"htmlUrl", "cloneUrl", "sshUrl", "fullName"} {
		if value, ok := scaffold.Repository[key]; ok {
			record = append(record, Field{Name: "Repository " + key, Value: JSONText(value)})
		}
	}
	for _, warning := range scaffold.Warnings {
		record = append(record, Field{Name: "Warning", Value: warning})
	}
	return record
}
