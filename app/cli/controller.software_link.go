package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SoftwareLinkControllerInterface is the "link" commands: the links of a
// software to elsewhere.
type SoftwareLinkControllerInterface interface {
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

// SoftwareLinkController is the [SoftwareLinkControllerInterface] over a
// [domain.SoftwareLinkServiceInterface].
type SoftwareLinkController struct {
	// SoftwareLinkService manages the links of the softwares.
	SoftwareLinkService domain.SoftwareLinkServiceInterface
}

// Index prints the links of SOFTWARE (an ID or a name).
func (c SoftwareLinkController) Index(cmd *cobra.Command, args []string) error {
	links, err := c.SoftwareLinkService.All(args[0])
	if err != nil {
		return err
	}
	return PrintList(cmd, links, SoftwareLinkSerializer{}.ModelToTable(links))
}

// Show prints the link ID.
func (c SoftwareLinkController) Show(cmd *cobra.Command, args []string) error {
	link, err := c.SoftwareLinkService.Show(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, SoftwareLinkSerializer{}.ModelToRecord(*link))
}

// Create adds a link from --label and --url to SOFTWARE and prints it.
func (c SoftwareLinkController) Create(cmd *cobra.Command, args []string) error {
	serializer := SoftwareLinkSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.SoftwareLinkService.Create(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Update changes the fields given as flags on the link ID and prints it.
func (c SoftwareLinkController) Update(cmd *cobra.Command, args []string) error {
	serializer := SoftwareLinkSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.SoftwareLinkService.Update(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Destroy deletes the link ID after asking (not with --yes).
func (c SoftwareLinkController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Delete the link "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.SoftwareLinkService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Link %s deleted.", args[0])
	return nil
}
