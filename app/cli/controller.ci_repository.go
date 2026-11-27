package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// CiRepositoryControllerInterface is the "ci-repo" commands: the CI
// repositories the engine knows and their runs.
type CiRepositoryControllerInterface interface {
	// Index lists the CI repositories.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one CI repository.
	Show(cmd *cobra.Command, args []string) error
	// Update sets a CI repository's platform or software.
	Update(cmd *cobra.Command, args []string) error
	// Sync reads the repositories from the CI server.
	Sync(cmd *cobra.Command, args []string) error
	// Runs lists the runs of a CI repository.
	Runs(cmd *cobra.Command, args []string) error
	// Trigger starts a run of a CI repository.
	Trigger(cmd *cobra.Command, args []string) error
}

// CiRepositoryController is the [CiRepositoryControllerInterface] over a
// [domain.CiRepositoryServiceInterface].
type CiRepositoryController struct {
	// CiRepositoryService manages the CI repositories.
	CiRepositoryService domain.CiRepositoryServiceInterface
}

// Index prints the CI repositories matching --platform and --active, one page
// or, with --all, every page.
func (c CiRepositoryController) Index(cmd *cobra.Command, args []string) error {
	serializer := CiRepositorySerializer{}
	page, err := c.CiRepositoryService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the CI repository ID.
func (c CiRepositoryController) Show(cmd *cobra.Command, args []string) error {
	ci_repository_id, err := domain.ParseID("CI repository", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	ci_repository, err := c.CiRepositoryService.Show(ci_repository_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, ci_repository, CiRepositorySerializer{}.ModelToRecord(*ci_repository))
}

// Update sets --platform and --software (an ID or name; "" detaches it)
// on the CI repository ID and prints it.
func (c CiRepositoryController) Update(cmd *cobra.Command, args []string) error {
	ci_repository_id, err := domain.ParseID("CI repository", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	serializer := CiRepositorySerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	ci_repository, err := c.CiRepositoryService.Update(ci_repository_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, ci_repository, serializer.ModelToRecord(*ci_repository))
}

// Sync syncs the CI repositories with the CI server and prints what changed.
func (c CiRepositoryController) Sync(cmd *cobra.Command, args []string) error {
	sync, err := c.CiRepositoryService.Sync()
	if err != nil {
		return err
	}
	return PrintList(cmd, sync, CiRepositorySerializer{}.SyncToTable(*sync))
}

// Runs prints one page (--page) of the runs of the CI repository ID, newest
// first.
func (c CiRepositoryController) Runs(cmd *cobra.Command, args []string) error {
	ci_repository_id, err := domain.ParseID("CI repository", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	page_number, _ := cmd.Flags().GetInt("page")
	page, err := c.CiRepositoryService.Runs(ci_repository_id, page_number)
	if err != nil {
		return err
	}
	return PrintList(cmd, page.Items, CiRepositorySerializer{}.RunsToTable(page.Items))
}

// Trigger starts a run of the CI repository ID on --branch and prints it.
func (c CiRepositoryController) Trigger(cmd *cobra.Command, args []string) error {
	ci_repository_id, err := domain.ParseID("CI repository", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	branch, _ := cmd.Flags().GetString("branch")
	run, err := c.CiRepositoryService.Trigger(ci_repository_id, branch)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, run, CiRepositorySerializer{}.RunToRecord(*run))
}
