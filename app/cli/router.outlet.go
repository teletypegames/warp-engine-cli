package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// outlet adds the "outlet" commands.
func (r Router) outlet(root *cobra.Command) {
	outlet := group(root, "operations", "outlet", "The outlets: the catalog API, shelves, Assembly64 …", "outlets")
	add(outlet, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the outlets, the catalog API first", Args: ExactArgs(0),
		Example: "  warp-engine outlet list\n  warp-engine outlet list -o json",
		Run:     r.OutletController.Index,
	})
	add(outlet, "", Route{
		Use: "check KEY", Short: "Show what an outlet offers and which softwares it leaves out", Args: ExactArgs(1, "KEY"),
		Example: "  warp-engine outlet check assembly64\n  warp-engine outlet check c64 --strict",
		Flags: func(flags *pflag.FlagSet) {
			flags.Bool("strict", false, "fail (exit code 1) when the outlet leaves a software out")
		},
		Run: r.OutletController.Check,
	})
}
