package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// download adds the "download" commands.
func (r Router) download(root *cobra.Command) {
	serializer := DownloadSerializer{}
	download := group(root, "operations", "download", "The download log and its statistics", "downloads")
	add(download, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the downloads, newest first", Args: ExactArgs(0),
		Example: "  warp-engine download list --software ra --from 2026-11-01\n  warp-engine download list --file-path .d64 --all -o json",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, DownloadSorts)
		},
		Run: r.DownloadController.Index,
	})
	add(download, "", Route{
		Use: "show ID", Short: "Show a download", Args: ExactArgs(1, "ID"),
		Run: r.DownloadController.Show,
	})
	add(download, "", Route{
		Use: "summary", Short: "Count the downloads: in total, today, this week, this month", Args: ExactArgs(0),
		Run: r.DownloadController.Summary,
	})
	add(download, "", Route{
		Use: "top-files", Short: "List the most downloaded files", Args: ExactArgs(0),
		Example: "  warp-engine download top-files --limit 20",
		Flags: func(flags *pflag.FlagSet) {
			flags.Int("limit", 0, "how many; the server's default is 10, its maximum 100")
		},
		Run: r.DownloadController.TopFiles,
	})
	add(download, "", Route{
		Use: "top-softwares", Short: "List the most downloaded softwares", Args: ExactArgs(0),
		Example: "  warp-engine download top-softwares --limit 10",
		Flags: func(flags *pflag.FlagSet) {
			flags.Int("limit", 0, "how many; the server's default is 5, its maximum 100")
		},
		Run: r.DownloadController.TopSoftwares,
	})
	add(download, "", Route{
		Use: "stats SOFTWARE", Short: "Show the downloads of a software (by ID or name), by release and file", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine download stats blessing-of-ra",
		Run:     r.DownloadController.Stats,
	})
}
