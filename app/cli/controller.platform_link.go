package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// PlatformLinkControllerInterface is the "platform-link" commands: the links shown with every software of a platform.
type PlatformLinkControllerInterface interface {
	// Index lists the platform links.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one platform link.
	Show(cmd *cobra.Command, args []string) error
	// Create creates a platform link.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a platform link.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a platform link.
	Destroy(cmd *cobra.Command, args []string) error
}

// PlatformLinkController is the [PlatformLinkControllerInterface] over a
// [domain.PlatformLinkServiceInterface].
type PlatformLinkController struct {
	// PlatformLinkService manages the platform links.
	PlatformLinkService domain.PlatformLinkServiceInterface
}

// Index prints the platform links matching the filter flags (--platform), one
// page or, with --all, every page.
func (c PlatformLinkController) Index(cmd *cobra.Command, args []string) error {
	serializer := PlatformLinkSerializer{}
	page, err := c.PlatformLinkService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the platform link ID.
func (c PlatformLinkController) Show(cmd *cobra.Command, args []string) error {
	platform_link_id, err := domain.ParseID("platform link", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	link, err := c.PlatformLinkService.Show(platform_link_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, PlatformLinkSerializer{}.ModelToRecord(*link))
}

// Create creates a platform link from the field flags (--platform, --name and --url are
// needed) and prints it.
func (c PlatformLinkController) Create(cmd *cobra.Command, args []string) error {
	serializer := PlatformLinkSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.PlatformLinkService.Create(attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Update changes the fields given as flags on the platform link ID and prints it.
func (c PlatformLinkController) Update(cmd *cobra.Command, args []string) error {
	platform_link_id, err := domain.ParseID("platform link", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	serializer := PlatformLinkSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	link, err := c.PlatformLinkService.Update(platform_link_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, link, serializer.ModelToRecord(*link))
}

// Destroy deletes the platform link ID after asking (not with --yes).
func (c PlatformLinkController) Destroy(cmd *cobra.Command, args []string) error {
	platform_link_id, err := domain.ParseID("platform link", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	ok, err := Confirm(cmd, "Delete the platform link "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.PlatformLinkService.Destroy(platform_link_id); err != nil {
		return err
	}
	PrintDone(cmd, "Platform link %s deleted.", args[0])
	return nil
}
