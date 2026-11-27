package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ciRepository adds the "ci-repo" commands.
func (r Router) ciRepository(root *cobra.Command) {
	serializer := CiRepositorySerializer{}
	ci_repository := group(root, "operations", "ci-repo", "The CI repositories the engine knows, and their pipeline runs", "ci-repos", "ci-repository")
	add(ci_repository, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the CI repositories", Args: ExactArgs(0),
		Example: "  warp-engine ci-repo list --platform c64 --active\n  warp-engine ci-repo list --sort -lastPipelineAt",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, CiRepositorySorts)
		},
		Run: r.CiRepositoryController.Index,
	})
	add(ci_repository, "", Route{
		Use: "show ID", Short: "Show a CI repository", Args: ExactArgs(1, "ID"),
		Run: r.CiRepositoryController.Show,
	})
	add(ci_repository, "", Route{
		Use: "update ID", Short: "Set a CI repository's platform or software", Args: ExactArgs(1, "ID"),
		Long: "Set the platform a CI repository builds for, or the software its builds are published as.\n" +
			"A software has one CI repository: giving it to this one takes it from the other.",
		Example: "  warp-engine ci-repo update 4 --software blessing-of-ra --platform c64\n  warp-engine ci-repo update 4 --software \"\"   # detaches the software",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.CiRepositoryController.Update,
	})
	add(ci_repository, "", Route{
		Use: "sync", Short: "Read the repositories from the CI server", Args: ExactArgs(0),
		Long: "Read the repositories from the CI server: new ones are added, changed ones\n" +
			"are updated, vanished ones deactivated. Prints what changed.",
		Run: r.CiRepositoryController.Sync,
	})
	add(ci_repository, "", Route{
		Use: "runs ID", Short: "List the pipeline runs of a CI repository, newest first", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine ci-repo runs 4\n  warp-engine ci-repo runs 4 --page 2",
		Flags:   func(flags *pflag.FlagSet) { flags.Int("page", 0, "the page of the CI server's listing, from 1") },
		Run:     r.CiRepositoryController.Runs,
	})
	add(ci_repository, "", Route{
		Use: "trigger ID", Short: "Start a pipeline run of a CI repository", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine ci-repo trigger 4\n  warp-engine ci-repo trigger 4 --branch develop",
		Flags: func(flags *pflag.FlagSet) {
			flags.String("branch", "", "the branch to build (default: the server's, main)")
		},
		Run: r.CiRepositoryController.Trigger,
	})
	if r.BuildPlatformController != nil {
		r.completeFlag(ci_repository, "platform", r.BuildPlatformController.CiRepositoryPlatforms)
	}
}
