package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SoftwareStatusControllerInterface is the "status" commands: the vocabulary of a software's status.
type SoftwareStatusControllerInterface interface {
	// Index lists the software statuses.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one status.
	Show(cmd *cobra.Command, args []string) error
	// Create creates a status.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a status.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a status.
	Destroy(cmd *cobra.Command, args []string) error
}

// SoftwareStatusController is the [SoftwareStatusControllerInterface] over a
// [domain.SoftwareStatusServiceInterface].
type SoftwareStatusController struct {
	// SoftwareStatusService manages the software statuses.
	SoftwareStatusService domain.SoftwareStatusServiceInterface
}

// Index prints the software statuses matching the filter flags (--default-visible), one
// page or, with --all, every page.
func (c SoftwareStatusController) Index(cmd *cobra.Command, args []string) error {
	serializer := SoftwareStatusSerializer{}
	page, err := c.SoftwareStatusService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the status ID.
func (c SoftwareStatusController) Show(cmd *cobra.Command, args []string) error {
	status_id, err := domain.ParseID("status", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	status, err := c.SoftwareStatusService.Show(status_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, status, SoftwareStatusSerializer{}.ModelToRecord(*status))
}

// Create creates a status from the field flags (--key and --label are
// needed) and prints it.
func (c SoftwareStatusController) Create(cmd *cobra.Command, args []string) error {
	serializer := SoftwareStatusSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	status, err := c.SoftwareStatusService.Create(attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, status, serializer.ModelToRecord(*status))
}

// Update changes the fields given as flags on the status ID and prints it.
func (c SoftwareStatusController) Update(cmd *cobra.Command, args []string) error {
	status_id, err := domain.ParseID("status", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	serializer := SoftwareStatusSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	status, err := c.SoftwareStatusService.Update(status_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, status, serializer.ModelToRecord(*status))
}

// Destroy deletes the status ID after asking (not with --yes).
func (c SoftwareStatusController) Destroy(cmd *cobra.Command, args []string) error {
	status_id, err := domain.ParseID("status", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	ok, err := Confirm(cmd, "Delete the status "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.SoftwareStatusService.Destroy(status_id); err != nil {
		return err
	}
	PrintDone(cmd, "Status %s deleted.", args[0])
	return nil
}
