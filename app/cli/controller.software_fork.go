package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SoftwareForkControllerInterface is the "fork" commands: the forks other
// engines announced of this engine's softwares, and their moderation.
type SoftwareForkControllerInterface interface {
	// Index lists the announced forks.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one announced fork.
	Show(cmd *cobra.Command, args []string) error
	// Update moderates an announced fork: sets its status.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes an announced fork.
	Destroy(cmd *cobra.Command, args []string) error
}

// SoftwareForkController is the [SoftwareForkControllerInterface] over a
// [domain.SoftwareForkServiceInterface].
type SoftwareForkController struct {
	// SoftwareForkService moderates the announced forks.
	SoftwareForkService domain.SoftwareForkServiceInterface
}

// Index prints the announced forks matching --status and --software, newest
// first, one page or, with --all, every page.
func (c SoftwareForkController) Index(cmd *cobra.Command, args []string) error {
	serializer := SoftwareForkSerializer{}
	page, err := c.SoftwareForkService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the announced fork ID.
func (c SoftwareForkController) Show(cmd *cobra.Command, args []string) error {
	software_fork_id, err := argID("fork", args[0])
	if err != nil {
		return err
	}
	link, err := c.SoftwareForkService.Show(software_fork_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, SoftwareForkSerializer{}.ModelToRecord(*link))
}

// Update sets the --status of the announced fork ID (visible lists it, hidden
// keeps it out) and prints it.
func (c SoftwareForkController) Update(cmd *cobra.Command, args []string) error {
	serializer := SoftwareForkSerializer{}
	software_fork_id, err := argID("fork", args[0])
	if err != nil {
		return err
	}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.SoftwareForkService.Update(software_fork_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Destroy deletes the announced fork ID after asking (not with --yes); a new
// announcement of the fork brings it back as pending.
func (c SoftwareForkController) Destroy(cmd *cobra.Command, args []string) error {
	software_fork_id, err := argID("fork", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Delete the announced fork "+args[0]+"? A new announcement brings it back as pending."); err != nil {
		return err
	}
	if err = c.SoftwareForkService.Destroy(software_fork_id); err != nil {
		return err
	}
	PrintDone(cmd, "Fork %d deleted.", software_fork_id)
	return nil
}
