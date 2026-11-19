package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SoftwareControllerInterface is the "software" commands: the entries of
// the catalog.
type SoftwareControllerInterface interface {
	// Index lists the softwares.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one software.
	Show(cmd *cobra.Command, args []string) error
	// Create creates a software.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a software.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a software.
	Destroy(cmd *cobra.Command, args []string) error
	// Scaffold sets a new software up with its repository and pipeline.
	Scaffold(cmd *cobra.Command, args []string) error
}

// SoftwareController is the [SoftwareControllerInterface] over a
// [domain.SoftwareServiceInterface].
type SoftwareController struct {
	// SoftwareService manages the softwares.
	SoftwareService domain.SoftwareServiceInterface
}

// Index prints the softwares matching the filter flags (--name, --platform,
// --query …), one page or, with --all, every page.
func (c SoftwareController) Index(cmd *cobra.Command, args []string) error {
	serializer := SoftwareSerializer{}
	page, err := c.SoftwareService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the software SOFTWARE (an ID or a name); --include embeds
// its releases, images, external links or extras.
func (c SoftwareController) Show(cmd *cobra.Command, args []string) error {
	include, _ := cmd.Flags().GetStringSlice("include")
	software, err := c.SoftwareService.Show(args[0], include)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, software, SoftwareSerializer{}.ModelToRecord(*software))
}

// Create creates a software from the field flags (--name, --title and
// --platform are needed) and prints it.
func (c SoftwareController) Create(cmd *cobra.Command, args []string) error {
	serializer := SoftwareSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	software, err := c.SoftwareService.Create(attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, software, serializer.ModelToRecord(*software))
}

// Update changes the fields given as flags on SOFTWARE and prints it.
func (c SoftwareController) Update(cmd *cobra.Command, args []string) error {
	serializer := SoftwareSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	software, err := c.SoftwareService.Update(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, software, serializer.ModelToRecord(*software))
}

// Destroy deletes SOFTWARE after asking (not with --yes).
func (c SoftwareController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Delete the software "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.SoftwareService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Software %s deleted.", args[0])
	return nil
}

// Scaffold creates a software with its forge repository and CI pipeline
// from --name, --title and --platform (and --category, --description) and
// prints what was made.
func (c SoftwareController) Scaffold(cmd *cobra.Command, args []string) error {
	serializer := SoftwareSerializer{}
	attributes, err := FieldsToAttributes(cmd.Flags(), serializer.ScaffoldFields())
	if err != nil {
		return err
	}
	scaffold, err := c.SoftwareService.Scaffold(attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, scaffold, serializer.ScaffoldToRecord(*scaffold))
}
