package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// softwareFork adds the "fork" commands.
func (r Router) softwareFork(root *cobra.Command) {
	serializer := SoftwareForkSerializer{}
	fork := group(root, "catalog", "fork", "The forks other engines announced of the softwares here", "forks")
	add(fork, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the announced forks, newest first", Args: ExactArgs(0),
		Example: "  warp-engine fork list --status pending\n  warp-engine fork list --software pong --sort -checkedAt",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, SoftwareForkSorts)
		},
		Run: r.SoftwareForkController.Index,
	})
	add(fork, "", Route{
		Use: "show ID", Short: "Show an announced fork", Args: ExactArgs(1, "ID"),
		Run: r.SoftwareForkController.Show,
	})
	add(fork, "", Route{
		Use: "update ID", Aliases: []string{"moderate"}, Short: "Moderate an announced fork: visible lists it, hidden keeps it out",
		Args:    ExactArgs(1, "ID"),
		Example: "  warp-engine fork update 4 --status visible\n  warp-engine fork moderate 5 --status hidden",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareForkController.Update,
	})
	add(fork, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete an announced fork; a new announcement brings it back as pending",
		Args: ExactArgs(1, "ID"), Flags: DeclareYes, Run: r.SoftwareForkController.Destroy,
	})
}
