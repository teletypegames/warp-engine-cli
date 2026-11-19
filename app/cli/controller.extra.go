package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// ExtraControllerInterface is the "extra" commands: the host's additional
// resources of a software, on hosted engines that declare some.
type ExtraControllerInterface interface {
	// Show shows what a resource holds for a software.
	Show(cmd *cobra.Command, args []string) error
	// Update replaces what a resource holds for a software.
	Update(cmd *cobra.Command, args []string) error
}

// ExtraController is the [ExtraControllerInterface] over a
// [domain.ExtraServiceInterface].
type ExtraController struct {
	// ExtraService manages the extras.
	ExtraService domain.ExtraServiceInterface
}

// Show prints what the resource ASSOCIATION holds for SOFTWARE (an ID or a name).
func (c ExtraController) Show(cmd *cobra.Command, args []string) error {
	extra, err := c.ExtraService.Show(args[0], args[1])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, extra, ExtraSerializer{}.ModelToRecord(*extra))
}

// Update replaces what the resource ASSOCIATION holds for SOFTWARE with
// --value (or the "value" of --data) and prints the result.
func (c ExtraController) Update(cmd *cobra.Command, args []string) error {
	serializer := ExtraSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	extra, err := c.ExtraService.Update(args[0], args[1], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, extra, serializer.ModelToRecord(*extra))
}
