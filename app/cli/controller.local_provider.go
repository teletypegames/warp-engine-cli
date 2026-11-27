package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// LocalProviderControllerInterface is the "local" commands: the
// providers of WarpEngine Local and the local CI of a software.
type LocalProviderControllerInterface interface {
	// Index lists the providers.
	Index(cmd *cobra.Command, args []string) error
	// Enable starts a provider.
	Enable(cmd *cobra.Command, args []string) error
	// Disable stops a provider.
	Disable(cmd *cobra.Command, args []string) error
	// Destroy removes a provider's containers.
	Destroy(cmd *cobra.Command, args []string) error
	// Credentials shows a provider's sign-in.
	Credentials(cmd *cobra.Command, args []string) error
	// SetupCI sets a software up on the local git server and CI.
	SetupCI(cmd *cobra.Command, args []string) error
	// TeardownCI stops the local CI of a software.
	TeardownCI(cmd *cobra.Command, args []string) error
}

// LocalProviderController is the [LocalProviderControllerInterface]
// over a [domain.LocalProviderServiceInterface]; it follows the
// tasks it starts with the [domain.LocalTaskServiceInterface].
type LocalProviderController struct {
	// LocalProviderService manages the providers.
	LocalProviderService domain.LocalProviderServiceInterface
	// LocalTaskService follows the tasks.
	LocalTaskService domain.LocalTaskServiceInterface
}

// Index prints the providers with their state.
func (c LocalProviderController) Index(cmd *cobra.Command, args []string) error {
	page, err := c.LocalProviderService.All()
	if err != nil {
		return err
	}
	return PrintList(cmd, page.Items, LocalProviderSerializer{}.ModelToTable(page.Items))
}

// Enable starts the provider KEY and waits for it (not with --no-wait).
func (c LocalProviderController) Enable(cmd *cobra.Command, args []string) error {
	task, err := c.LocalProviderService.Enable(args[0])
	if err != nil {
		return err
	}
	return c.follow(cmd, task)
}

// Disable stops the provider KEY and waits for it (not with --no-wait).
func (c LocalProviderController) Disable(cmd *cobra.Command, args []string) error {
	task, err := c.LocalProviderService.Disable(args[0])
	if err != nil {
		return err
	}
	return c.follow(cmd, task)
}

// Destroy removes the containers of the provider KEY (with --purge its
// data too) after asking (not with --yes), and waits for it (not with
// --no-wait).
func (c LocalProviderController) Destroy(cmd *cobra.Command, args []string) error {
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
	task, err := c.LocalProviderService.Remove(args[0], purge)
	if err != nil {
		return err
	}
	return c.follow(cmd, task)
}

// Credentials prints the generated sign-in of the provider KEY.
func (c LocalProviderController) Credentials(cmd *cobra.Command, args []string) error {
	credentials, err := c.LocalProviderService.Credentials(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, credentials, LocalProviderSerializer{}.CredentialsToRecord(*credentials))
}

// SetupCI sets SOFTWARE (an ID or name) up on the local git server and CI
// and waits for it (not with --no-wait); the result holds the push URL.
func (c LocalProviderController) SetupCI(cmd *cobra.Command, args []string) error {
	task, err := c.LocalProviderService.SetupCI(args[0])
	if err != nil {
		return err
	}
	return c.follow(cmd, task)
}

// TeardownCI stops the local CI of SOFTWARE after asking (not with --yes);
// its repository stays.
func (c LocalProviderController) TeardownCI(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Stop the local CI of "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.LocalProviderService.TeardownCI(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Local CI of %s stopped; its repository stays.", args[0])
	return nil
}

// follow prints task, after waiting for it unless --no-wait is set.
func (c LocalProviderController) follow(cmd *cobra.Command, task *domain.LocalTaskModel) error {
	no_wait, _ := cmd.Flags().GetBool("no-wait")
	return followTask(cmd, c.LocalTaskService, task, !no_wait)
}
