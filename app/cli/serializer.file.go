package cli

import "git.teletypegames.org/engines/warp-engine-cli/app/domain"

// FileSerializer turns the entries of the artifact store into tables and
// records.
type FileSerializer struct{}

// ModelToTable renders directory entries as a table: path, type, size,
// modified, downloads. A directory's path ends in "/".
func (s FileSerializer) ModelToTable(entries []domain.FileEntryModel) Table {
	table := Table{Headers: []string{"Path", "Type", "Size", "Modified", "Downloads"}}
	for _, entry := range entries {
		entry_path := entry.Path
		size := "-"
		if entry.Type == "directory" {
			entry_path += "/"
		} else if entry.Size != nil {
			size = Bytes(*entry.Size)
		}
		table.Rows = append(table.Rows, []string{entry_path, entry.Type, size, Text(entry.ModifiedAt), Number(entry.Downloads)})
	}
	return table
}

// PathToRecord renders the answer of a change in the store: its path.
func (s FileSerializer) PathToRecord(file_path domain.FilePathModel) Record {
	return Record{{Name: "Path", Value: file_path.Path}}
}
