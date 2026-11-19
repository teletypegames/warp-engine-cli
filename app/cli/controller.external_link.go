package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// ExternalLinkControllerInterface is the "link" commands: the links of a
// software to elsewhere.
type ExternalLinkControllerInterface interface {
	// Index lists the links of a software.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one link.
	Show(cmd *cobra.Command, args []string) error
	// Create adds a link to a software.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a link.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a link.
	Destroy(cmd *cobra.Command, args []string) error
}

// ExternalLinkController is the [ExternalLinkControllerInterface] over a
// [domain.ExternalLinkServiceInterface].
type ExternalLinkController struct {
	// ExternalLinkService manages the external links.
	ExternalLinkService domain.ExternalLinkServiceInterface
}

// Index prints the links of SOFTWARE (an ID or a name).
func (c ExternalLinkController) Index(cmd *cobra.Command, args []string) error {
	links, err := c.ExternalLinkService.All(args[0])
	if err != nil {
		return err
	}
	return PrintList(cmd, links, ExternalLinkSerializer{}.ModelToTable(links))
}

// Show prints the link ID.
func (c ExternalLinkController) Show(cmd *cobra.Command, args []string) error {
	link, err := c.ExternalLinkService.Show(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, ExternalLinkSerializer{}.ModelToRecord(*link))
}

// Create adds a link from --label and --url to SOFTWARE and prints it.
func (c ExternalLinkController) Create(cmd *cobra.Command, args []string) error {
	serializer := ExternalLinkSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.ExternalLinkService.Create(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Update changes the fields given as flags on the link ID and prints it.
func (c ExternalLinkController) Update(cmd *cobra.Command, args []string) error {
	serializer := ExternalLinkSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.ExternalLinkService.Update(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Destroy deletes the link ID after asking (not with --yes).
func (c ExternalLinkController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Delete the link "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.ExternalLinkService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Link %s deleted.", args[0])
	return nil
}
