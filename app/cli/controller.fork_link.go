package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// ForkLinkControllerInterface is the "fork" commands: the forks other
// engines announced of this engine's softwares, and their moderation.
type ForkLinkControllerInterface interface {
	// Index lists the fork links.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one fork link.
	Show(cmd *cobra.Command, args []string) error
	// Update moderates a fork link: sets its status.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a fork link.
	Destroy(cmd *cobra.Command, args []string) error
}

// ForkLinkController is the [ForkLinkControllerInterface] over a
// [domain.ForkLinkServiceInterface].
type ForkLinkController struct {
	// ForkLinkService moderates the fork links.
	ForkLinkService domain.ForkLinkServiceInterface
}

// Index prints the fork links matching --status and --software, newest
// first, one page or, with --all, every page.
func (c ForkLinkController) Index(cmd *cobra.Command, args []string) error {
	serializer := ForkLinkSerializer{}
	page, err := c.ForkLinkService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the fork link ID.
func (c ForkLinkController) Show(cmd *cobra.Command, args []string) error {
	fork_link_id, err := argID("fork link", args[0])
	if err != nil {
		return err
	}
	link, err := c.ForkLinkService.Show(fork_link_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, ForkLinkSerializer{}.ModelToRecord(*link))
}

// Update sets the --status of the fork link ID (visible lists it, hidden
// keeps it out) and prints it.
func (c ForkLinkController) Update(cmd *cobra.Command, args []string) error {
	serializer := ForkLinkSerializer{}
	fork_link_id, err := argID("fork link", args[0])
	if err != nil {
		return err
	}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.ForkLinkService.Update(fork_link_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Destroy deletes the fork link ID after asking (not with --yes); a new
// announcement of the fork brings it back as pending.
func (c ForkLinkController) Destroy(cmd *cobra.Command, args []string) error {
	fork_link_id, err := argID("fork link", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Delete the fork link "+args[0]+"? A new announcement brings it back as pending."); err != nil {
		return err
	}
	if err = c.ForkLinkService.Destroy(fork_link_id); err != nil {
		return err
	}
	PrintDone(cmd, "Fork link %d deleted.", fork_link_id)
	return nil
}
