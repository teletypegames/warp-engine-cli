package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// runtime adds the "runtime" commands.
func (r Router) runtime(root *cobra.Command) {
	serializer := RuntimeSerializer{}
	runtime := group(root, "operations", "runtime", "The runtimes games are played in on the web", "runtimes")
	add(runtime, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the installed runtime versions", Args: ExactArgs(0),
		Example: "  warp-engine runtime list --name tic80",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, RuntimeSorts)
		},
		Run: r.RuntimeController.Index,
	})
	add(runtime, "", Route{
		Use: "show ID", Short: "Show an installed runtime version", Args: ExactArgs(1, "ID"),
		Run: r.RuntimeController.Show,
	})
	add(runtime, "", Route{
		Use: "install NAME [VERSION]", Short: "Install a downloaded runtime and the components it needs", Args: RangeArgs(1, 2, "NAME", "[VERSION]"),
		Long: "Install a runtime the engine downloads (not an uploaded one) and the components\n" +
			"its platforms need. Without VERSION the engine installs the version it has chosen.\n" +
			"With --dry-run nothing is installed: it prints what would be fetched.",
		Example: "  warp-engine runtime install tic80 --dry-run\n  warp-engine runtime install tic80 1.2.2715-dev",
		Flags:   func(flags *pflag.FlagSet) { flags.Bool("dry-run", false, "only report what would be fetched") },
		Run:     r.RuntimeController.Install,
	})
}
