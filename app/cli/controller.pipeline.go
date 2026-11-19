package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// PipelineControllerInterface is the "pipeline" commands: the CI
// repositories the engine knows and their runs.
type PipelineControllerInterface interface {
	// Index lists the pipelines.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one pipeline.
	Show(cmd *cobra.Command, args []string) error
	// Update sets a pipeline's platform or software.
	Update(cmd *cobra.Command, args []string) error
	// Sync reads the repositories from the CI server.
	Sync(cmd *cobra.Command, args []string) error
	// Runs lists the runs of a pipeline.
	Runs(cmd *cobra.Command, args []string) error
	// Trigger starts a run of a pipeline.
	Trigger(cmd *cobra.Command, args []string) error
}

// PipelineController is the [PipelineControllerInterface] over a
// [domain.PipelineServiceInterface].
type PipelineController struct {
	// PipelineService manages the pipelines.
	PipelineService domain.PipelineServiceInterface
}

// Index prints the pipelines matching --platform and --active, one page
// or, with --all, every page.
func (c PipelineController) Index(cmd *cobra.Command, args []string) error {
	serializer := PipelineSerializer{}
	page, err := c.PipelineService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the pipeline ID.
func (c PipelineController) Show(cmd *cobra.Command, args []string) error {
	pipeline_id, err := domain.ParseID("pipeline", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	pipeline, err := c.PipelineService.Show(pipeline_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, pipeline, PipelineSerializer{}.ModelToRecord(*pipeline))
}

// Update sets --platform and --software (an ID or name; "" detaches it)
// on the pipeline ID and prints it.
func (c PipelineController) Update(cmd *cobra.Command, args []string) error {
	pipeline_id, err := domain.ParseID("pipeline", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	serializer := PipelineSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	pipeline, err := c.PipelineService.Update(pipeline_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, pipeline, serializer.ModelToRecord(*pipeline))
}

// Sync syncs the pipelines with the CI server and prints what changed.
func (c PipelineController) Sync(cmd *cobra.Command, args []string) error {
	sync, err := c.PipelineService.Sync()
	if err != nil {
		return err
	}
	return PrintList(cmd, sync, PipelineSerializer{}.SyncToTable(*sync))
}

// Runs prints one page (--page) of the runs of the pipeline ID, newest
// first.
func (c PipelineController) Runs(cmd *cobra.Command, args []string) error {
	pipeline_id, err := domain.ParseID("pipeline", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	page_number, _ := cmd.Flags().GetInt("page")
	page, err := c.PipelineService.Runs(pipeline_id, page_number)
	if err != nil {
		return err
	}
	return PrintList(cmd, page.Items, PipelineSerializer{}.RunsToTable(page.Items))
}

// Trigger starts a run of the pipeline ID on --branch and prints it.
func (c PipelineController) Trigger(cmd *cobra.Command, args []string) error {
	pipeline_id, err := domain.ParseID("pipeline", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	branch, _ := cmd.Flags().GetString("branch")
	run, err := c.PipelineService.Trigger(pipeline_id, branch)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, run, PipelineSerializer{}.RunToRecord(*run))
}
