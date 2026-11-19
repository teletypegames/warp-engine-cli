package cli

import (
	"fmt"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// LocalOperationControllerInterface is "local operation": the background
// jobs of WarpEngine Local.
type LocalOperationControllerInterface interface {
	// Show shows one operation, or follows it with --wait.
	Show(cmd *cobra.Command, args []string) error
}

// LocalOperationController is the [LocalOperationControllerInterface] over
// a [domain.LocalOperationServiceInterface].
type LocalOperationController struct {
	// LocalOperationService follows the operations.
	LocalOperationService domain.LocalOperationServiceInterface
}

// Show prints the operation ID; with --wait it first follows it until it
// has finished, like the commands that start one.
func (c LocalOperationController) Show(cmd *cobra.Command, args []string) error {
	operation_id, err := domain.ParseID("operation", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	operation, err := c.LocalOperationService.Show(operation_id)
	if err != nil {
		return err
	}
	wait, _ := cmd.Flags().GetBool("wait")
	return followOperation(cmd, c.LocalOperationService, operation, wait)
}

// followOperation prints operation; with wait it first polls it until it
// has finished (at most --timeout), printing every new step on standard
// error as it comes. A failed or unfinished operation is printed and then
// returned as the error.
func followOperation(cmd *cobra.Command, service domain.LocalOperationServiceInterface, operation *domain.LocalOperationModel, wait bool) error {
	serializer := LocalOperationSerializer{}
	if !wait {
		return PrintRecord(cmd, operation, serializer.ModelToRecord(*operation))
	}
	timeout, _ := cmd.Flags().GetDuration("timeout")
	printed := 0
	progress := func(current domain.LocalOperationModel) {
		if printed > len(current.Steps) {
			printed = 0
		}
		for _, step := range current.Steps[printed:] {
			fmt.Fprintln(cmd.ErrOrStderr(), serializer.StepText(step))
		}
		printed = len(current.Steps)
	}
	final, err := service.Wait(*operation, timeout, progress)
	if final != nil {
		if print_err := PrintRecord(cmd, final, serializer.ModelToRecord(*final)); print_err != nil && err == nil {
			err = print_err
		}
	}
	return err
}
