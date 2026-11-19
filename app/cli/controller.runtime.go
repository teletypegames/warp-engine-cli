package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// RuntimeControllerInterface is the "runtime" commands: the installed
// runtimes.
type RuntimeControllerInterface interface {
	// Index lists the runtime versions.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one runtime version.
	Show(cmd *cobra.Command, args []string) error
	// Install installs a downloaded runtime.
	Install(cmd *cobra.Command, args []string) error
}

// RuntimeController is the [RuntimeControllerInterface] over a
// [domain.RuntimeServiceInterface].
type RuntimeController struct {
	// RuntimeService manages the runtimes.
	RuntimeService domain.RuntimeServiceInterface
}

// Index prints the runtime versions matching --name, one page or, with
// --all, every page.
func (c RuntimeController) Index(cmd *cobra.Command, args []string) error {
	serializer := RuntimeSerializer{}
	page, err := c.RuntimeService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the runtime version ID.
func (c RuntimeController) Show(cmd *cobra.Command, args []string) error {
	runtime_id, err := domain.ParseID("runtime", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	runtime, err := c.RuntimeService.Show(runtime_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, runtime, RuntimeSerializer{}.ModelToRecord(*runtime))
}

// Install installs NAME, at VERSION when given, or with --dry-run only
// reports what it would fetch, and prints the result.
func (c RuntimeController) Install(cmd *cobra.Command, args []string) error {
	version := ""
	if len(args) > 1 {
		version = args[1]
	}
	dry_run, _ := cmd.Flags().GetBool("dry-run")
	install, err := c.RuntimeService.Install(args[0], version, dry_run)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, install, RuntimeSerializer{}.InstallToRecord(*install))
}
