package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// ReleaseControllerInterface is the "release" commands: the versions of a
// software.
type ReleaseControllerInterface interface {
	// Index lists the releases of a software.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one release with its files.
	Show(cmd *cobra.Command, args []string) error
	// Create adds a release to a software.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a release.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a release with its files.
	Destroy(cmd *cobra.Command, args []string) error
}

// ReleaseController is the [ReleaseControllerInterface] over a
// [domain.ReleaseServiceInterface].
type ReleaseController struct {
	// ReleaseService manages the releases.
	ReleaseService domain.ReleaseServiceInterface
}

// Index prints the releases of SOFTWARE (an ID or a name) matching
// --channel and --commit-sha, one page or, with --all, every page.
func (c ReleaseController) Index(cmd *cobra.Command, args []string) error {
	serializer := ReleaseSerializer{}
	page, err := c.ReleaseService.All(args[0], serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the release ID with its assets.
func (c ReleaseController) Show(cmd *cobra.Command, args []string) error {
	release, err := c.ReleaseService.Show(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, release, ReleaseSerializer{}.ModelToRecord(*release))
}

// Create adds a release with --version to SOFTWARE and prints it.
func (c ReleaseController) Create(cmd *cobra.Command, args []string) error {
	serializer := ReleaseSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	release, err := c.ReleaseService.Create(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, release, serializer.ModelToRecord(*release))
}

// Update changes the fields given as flags on the release ID and prints it.
func (c ReleaseController) Update(cmd *cobra.Command, args []string) error {
	serializer := ReleaseSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	release, err := c.ReleaseService.Update(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, release, serializer.ModelToRecord(*release))
}

// Destroy deletes the release ID with its assets after asking (not with --yes).
func (c ReleaseController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Delete the release "+args[0]+" with its files?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.ReleaseService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Release %s deleted.", args[0])
	return nil
}
