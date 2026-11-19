package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// DownloadControllerInterface is the "download" commands: the download log
// and its statistics.
type DownloadControllerInterface interface {
	// Index lists the downloads.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one download.
	Show(cmd *cobra.Command, args []string) error
	// Summary shows the download counts.
	Summary(cmd *cobra.Command, args []string) error
	// TopFiles lists the most downloaded files.
	TopFiles(cmd *cobra.Command, args []string) error
	// TopSoftwares lists the most downloaded softwares.
	TopSoftwares(cmd *cobra.Command, args []string) error
	// Stats shows the downloads of one software.
	Stats(cmd *cobra.Command, args []string) error
}

// DownloadController is the [DownloadControllerInterface] over a
// [domain.DownloadServiceInterface].
type DownloadController struct {
	// DownloadService reads the download log.
	DownloadService domain.DownloadServiceInterface
}

// Index prints the downloads matching the filter flags (--file-path,
// --software, --from, --to …), newest first, one page or, with --all,
// every page.
func (c DownloadController) Index(cmd *cobra.Command, args []string) error {
	serializer := DownloadSerializer{}
	page, err := c.DownloadService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the download ID.
func (c DownloadController) Show(cmd *cobra.Command, args []string) error {
	download_id, err := domain.ParseID("download", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	download, err := c.DownloadService.Show(download_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, download, DownloadSerializer{}.ModelToRecord(*download))
}

// Summary prints the downloads in total, today, this week and this month.
func (c DownloadController) Summary(cmd *cobra.Command, args []string) error {
	summary, err := c.DownloadService.Summary()
	if err != nil {
		return err
	}
	return PrintRecord(cmd, summary, DownloadSerializer{}.SummaryToRecord(*summary))
}

// TopFiles prints the --limit most downloaded files.
func (c DownloadController) TopFiles(cmd *cobra.Command, args []string) error {
	limit, _ := cmd.Flags().GetInt("limit")
	files, err := c.DownloadService.TopFiles(limit)
	if err != nil {
		return err
	}
	return PrintList(cmd, files, DownloadSerializer{}.FilesToTable(files))
}

// TopSoftwares prints the --limit most downloaded softwares.
func (c DownloadController) TopSoftwares(cmd *cobra.Command, args []string) error {
	limit, _ := cmd.Flags().GetInt("limit")
	softwares, err := c.DownloadService.TopSoftwares(limit)
	if err != nil {
		return err
	}
	return PrintList(cmd, softwares, DownloadSerializer{}.SoftwaresToTable(softwares))
}

// Stats prints the downloads of SOFTWARE (an ID or a name) by release and
// by file.
func (c DownloadController) Stats(cmd *cobra.Command, args []string) error {
	stats, err := c.DownloadService.Stats(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, stats, DownloadSerializer{}.StatsToRecord(*stats))
}
