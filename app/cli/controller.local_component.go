package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// LocalComponentControllerInterface is the "local" commands: the
// components of WarpEngine Local and the local CI of a software.
type LocalComponentControllerInterface interface {
	// Index lists the components.
	Index(cmd *cobra.Command, args []string) error
	// Enable starts a component.
	Enable(cmd *cobra.Command, args []string) error
	// Disable stops a component.
	Disable(cmd *cobra.Command, args []string) error
	// Destroy removes a component's containers.
	Destroy(cmd *cobra.Command, args []string) error
	// Credentials shows a component's sign-in.
	Credentials(cmd *cobra.Command, args []string) error
	// SetupCI sets a software up on the local git server and CI.
	SetupCI(cmd *cobra.Command, args []string) error
	// TeardownCI stops the local CI of a software.
	TeardownCI(cmd *cobra.Command, args []string) error
}

// LocalComponentController is the [LocalComponentControllerInterface]
// over a [domain.LocalComponentServiceInterface]; it follows the
// operations it starts with the [domain.LocalOperationServiceInterface].
type LocalComponentController struct {
	// LocalComponentService manages the components.
	LocalComponentService domain.LocalComponentServiceInterface
	// LocalOperationService follows the operations.
	LocalOperationService domain.LocalOperationServiceInterface
}

// Index prints the components with their state.
func (c LocalComponentController) Index(cmd *cobra.Command, args []string) error {
	page, err := c.LocalComponentService.All()
	if err != nil {
		return err
	}
	return PrintList(cmd, page.Items, LocalComponentSerializer{}.ModelToTable(page.Items))
}

// Enable starts the component KEY and waits for it (not with --no-wait).
func (c LocalComponentController) Enable(cmd *cobra.Command, args []string) error {
	operation, err := c.LocalComponentService.Enable(args[0])
	if err != nil {
		return err
	}
	return c.follow(cmd, operation)
}

// Disable stops the component KEY and waits for it (not with --no-wait).
func (c LocalComponentController) Disable(cmd *cobra.Command, args []string) error {
	operation, err := c.LocalComponentService.Disable(args[0])
	if err != nil {
		return err
	}
	return c.follow(cmd, operation)
}

// Destroy removes the containers of the component KEY (with --purge its
// data too) after asking (not with --yes), and waits for it (not with
// --no-wait).
func (c LocalComponentController) Destroy(cmd *cobra.Command, args []string) error {
	purge, _ := cmd.Flags().GetBool("purge")
	question := "Remove the containers of " + args[0] + "?"
	if purge {
		question = "Remove " + args[0] + " with its data and sign-ins?"
	}
	ok, err := Confirm(cmd, question)
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	operation, err := c.LocalComponentService.Remove(args[0], purge)
	if err != nil {
		return err
	}
	return c.follow(cmd, operation)
}

// Credentials prints the generated sign-in of the component KEY.
func (c LocalComponentController) Credentials(cmd *cobra.Command, args []string) error {
	credentials, err := c.LocalComponentService.Credentials(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, credentials, LocalComponentSerializer{}.CredentialsToRecord(*credentials))
}

// SetupCI sets SOFTWARE (an ID or name) up on the local git server and CI
// and waits for it (not with --no-wait); the result holds the push URL.
func (c LocalComponentController) SetupCI(cmd *cobra.Command, args []string) error {
	operation, err := c.LocalComponentService.SetupCI(args[0])
	if err != nil {
		return err
	}
	return c.follow(cmd, operation)
}

// TeardownCI stops the local CI of SOFTWARE after asking (not with --yes);
// its repository stays.
func (c LocalComponentController) TeardownCI(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Stop the local CI of "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.LocalComponentService.TeardownCI(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Local CI of %s stopped; its repository stays.", args[0])
	return nil
}

// follow prints operation, after waiting for it unless --no-wait is set.
func (c LocalComponentController) follow(cmd *cobra.Command, operation *domain.LocalOperationModel) error {
	no_wait, _ := cmd.Flags().GetBool("no-wait")
	return followOperation(cmd, c.LocalOperationService, operation, !no_wait)
}
