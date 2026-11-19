package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// FileControllerInterface is the "file" commands: the artifact store.
type FileControllerInterface interface {
	// Index lists a directory.
	Index(cmd *cobra.Command, args []string) error
	// Upload uploads a local file.
	Upload(cmd *cobra.Command, args []string) error
	// Rename renames a file or directory in place.
	Rename(cmd *cobra.Command, args []string) error
	// Destroy deletes a file or an empty directory.
	Destroy(cmd *cobra.Command, args []string) error
	// CreateDirectory makes a directory.
	CreateDirectory(cmd *cobra.Command, args []string) error
}

// FileController is the [FileControllerInterface] over a
// [domain.FileServiceInterface].
type FileController struct {
	// FileService manages the artifact store.
	FileService domain.FileServiceInterface
}

// Index prints the directory DIR (the root without it).
func (c FileController) Index(cmd *cobra.Command, args []string) error {
	dir := ""
	if len(args) > 0 {
		dir = args[0]
	}
	page, err := c.FileService.All(dir)
	if err != nil {
		return err
	}
	return PrintList(cmd, page.Items, FileSerializer{}.ModelToTable(page.Items))
}

// Upload sends the local FILE into --dir and prints its path in the store.
func (c FileController) Upload(cmd *cobra.Command, args []string) error {
	dir, _ := cmd.Flags().GetString("dir")
	file_path, err := c.FileService.Upload(args[0], dir)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, file_path, FileSerializer{}.PathToRecord(*file_path))
}

// Rename gives PATH the name NAME and prints the new path.
func (c FileController) Rename(cmd *cobra.Command, args []string) error {
	file_path, err := c.FileService.Rename(args[0], args[1])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, file_path, FileSerializer{}.PathToRecord(*file_path))
}

// Destroy deletes PATH after asking (not with --yes).
func (c FileController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Delete "+args[0]+" from the artifact store?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.FileService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "%s deleted.", args[0])
	return nil
}

// CreateDirectory makes the directory PATH and prints its path.
func (c FileController) CreateDirectory(cmd *cobra.Command, args []string) error {
	file_path, err := c.FileService.CreateDirectory(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, file_path, FileSerializer{}.PathToRecord(*file_path))
}
