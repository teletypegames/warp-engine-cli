package cli

import (
	"fmt"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// LocalTaskControllerInterface is "local task": the background
// jobs of WarpEngine Local.
type LocalTaskControllerInterface interface {
	// Show shows one task, or follows it with --wait.
	Show(cmd *cobra.Command, args []string) error
}

// LocalTaskController is the [LocalTaskControllerInterface] over
// a [domain.LocalTaskServiceInterface].
type LocalTaskController struct {
	// LocalTaskService follows the tasks.
	LocalTaskService domain.LocalTaskServiceInterface
}

// Show prints the task ID; with --wait it first follows it until it
// has finished, like the commands that start one.
func (c LocalTaskController) Show(cmd *cobra.Command, args []string) error {
	task_id, err := domain.ParseID("task", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	task, err := c.LocalTaskService.Show(task_id)
	if err != nil {
		return err
	}
	wait, _ := cmd.Flags().GetBool("wait")
	return followTask(cmd, c.LocalTaskService, task, wait)
}

// followTask prints task; with wait it first polls it until it
// has finished (at most --timeout), printing every new step on standard
// error as it comes. A failed or unfinished task is printed and then
// returned as the error.
func followTask(cmd *cobra.Command, service domain.LocalTaskServiceInterface, task *domain.LocalTaskModel, wait bool) error {
	serializer := LocalTaskSerializer{}
	if !wait {
		return PrintRecord(cmd, task, serializer.ModelToRecord(*task))
	}
	timeout, _ := cmd.Flags().GetDuration("timeout")
	printed := 0
	progress := func(current domain.LocalTaskModel) {
		if printed > len(current.Steps) {
			printed = 0
		}
		for _, step := range current.Steps[printed:] {
			fmt.Fprintln(cmd.ErrOrStderr(), serializer.StepText(step))
		}
		printed = len(current.Steps)
	}
	final, err := service.Wait(*task, timeout, progress)
	if final != nil {
		if print_err := PrintRecord(cmd, final, serializer.ModelToRecord(*final)); print_err != nil && err == nil {
			err = print_err
		}
	}
	return err
}
