package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// localTask adds "task" under the "local" command local.
func (r Router) localTask(local *cobra.Command) {
	add(local, "", Route{
		Use: "task ID", Short: "Show a background task, or follow it with --wait", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine local task 12 --wait",
		Flags: func(flags *pflag.FlagSet) {
			flags.Bool("wait", false, "follow it until it has finished, printing its steps")
			flags.Duration("timeout", DefaultTaskTimeout, "how long to wait with --wait")
		},
		Run: r.LocalTaskController.Show,
	})
}
