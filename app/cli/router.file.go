package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// file adds the "file" commands.
func (r Router) file(root *cobra.Command) {
	file := group(root, "operations", "file", "The artifact store: the files the engine serves", "files")
	add(file, "", Route{
		Use: "list [DIR]", Aliases: []string{"ls"}, Short: "List a directory of the artifact store", Args: RangeArgs(0, 1, "[DIR]"),
		Example: "  warp-engine file list\n  warp-engine file list blessing-of-ra/1.2.0",
		Run:     r.FileController.Index,
	})
	add(file, "", Route{
		Use: "upload FILE", Short: "Upload a local file into a directory", Args: ExactArgs(1, "FILE"),
		Example: "  warp-engine file upload ./manual.pdf --dir blessing-of-ra/docs",
		Flags: func(flags *pflag.FlagSet) {
			flags.String("dir", "", "the directory to upload into (default: the root)")
		},
		Run: r.FileController.Upload,
	})
	add(file, "", Route{
		Use: "rename PATH NAME", Short: "Rename a file or directory in its directory", Args: ExactArgs(2, "PATH", "NAME"),
		Long: "Rename a file or directory in place: NAME is the new name, not a path; the\n" +
			"artifact store does not move entries between directories.",
		Example: "  warp-engine file rename docs/manual.pdf handbook.pdf",
		Run:     r.FileController.Rename,
	})
	add(file, "", Route{
		Use: "delete PATH", Aliases: []string{"rm"}, Short: "Delete a file or an empty directory", Args: ExactArgs(1, "PATH"),
		Flags: DeclareYes, Run: r.FileController.Destroy,
	})
	add(file, "", Route{
		Use: "mkdir PATH", Short: "Create a directory (its parent must exist)", Args: ExactArgs(1, "PATH"),
		Example: "  warp-engine file mkdir blessing-of-ra/docs",
		Run:     r.FileController.CreateDirectory,
	})
}
