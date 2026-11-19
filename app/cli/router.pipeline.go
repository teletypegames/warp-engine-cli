package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// pipeline adds the "pipeline" commands.
func (r Router) pipeline(root *cobra.Command) {
	serializer := PipelineSerializer{}
	pipeline := group(root, "operations", "pipeline", "The CI pipelines the engine knows, and their runs", "pipelines")
	add(pipeline, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the pipelines", Args: ExactArgs(0),
		Example: "  warp-engine pipeline list --platform c64 --active\n  warp-engine pipeline list --sort -lastPipelineAt",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, PipelineSorts)
		},
		Run: r.PipelineController.Index,
	})
	add(pipeline, "", Route{
		Use: "show ID", Short: "Show a pipeline", Args: ExactArgs(1, "ID"),
		Run: r.PipelineController.Show,
	})
	add(pipeline, "", Route{
		Use: "update ID", Short: "Set a pipeline's platform or software", Args: ExactArgs(1, "ID"),
		Long: "Set the platform a pipeline builds for, or the software its builds are published as.\n" +
			"A software has one pipeline: giving it to this one takes it from the other.",
		Example: "  warp-engine pipeline update 4 --software blessing-of-ra --platform c64\n  warp-engine pipeline update 4 --software \"\"   # detaches the software",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.PipelineController.Update,
	})
	add(pipeline, "", Route{
		Use: "sync", Short: "Read the repositories from the CI server", Args: ExactArgs(0),
		Long: "Read the repositories from the CI server: new ones become pipelines, changed ones\n" +
			"are updated, vanished ones deactivated. Prints what changed.",
		Run: r.PipelineController.Sync,
	})
	add(pipeline, "", Route{
		Use: "runs ID", Short: "List the runs of a pipeline, newest first", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine pipeline runs 4\n  warp-engine pipeline runs 4 --page 2",
		Flags:   func(flags *pflag.FlagSet) { flags.Int("page", 0, "the page of the CI server's listing, from 1") },
		Run:     r.PipelineController.Runs,
	})
	add(pipeline, "", Route{
		Use: "trigger ID", Short: "Start a run of a pipeline", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine pipeline trigger 4\n  warp-engine pipeline trigger 4 --branch develop",
		Flags: func(flags *pflag.FlagSet) {
			flags.String("branch", "", "the branch to build (default: the server's, main)")
		},
		Run: r.PipelineController.Trigger,
	})
}
