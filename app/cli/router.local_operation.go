package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// localOperation adds "operation" under the "local" command local.
func (r Router) localOperation(local *cobra.Command) {
	add(local, "", Route{
		Use: "operation ID", Short: "Show a background operation, or follow it with --wait", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine local operation 12 --wait",
		Flags: func(flags *pflag.FlagSet) {
			flags.Bool("wait", false, "follow it until it has finished, printing its steps")
			flags.Duration("timeout", DefaultOperationTimeout, "how long to wait with --wait")
		},
		Run: r.LocalOperationController.Show,
	})
}
