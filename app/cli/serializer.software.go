package cli

import (
	"sort"
	"strconv"

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
	}
}

// ScaffoldFields are the field flags of "software scaffold".
func (s SoftwareSerializer) ScaffoldFields() []FieldFlag {
	return []FieldFlag{
		{Flag: "name", Key: "name", Usage: "software and repository name: a-z, 0-9 and dashes"},
		{Flag: "title", Key: "title", Usage: "display name"},
		{Flag: "platform", Key: "platform", Usage: "platform", Enum: domain.Platforms},
		{Flag: "category", Key: "category", Kind: KindNullable, Usage: "a category key"},
		{Flag: "description", Key: "description", Kind: KindNullable, Usage: "short description"},
	}
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
