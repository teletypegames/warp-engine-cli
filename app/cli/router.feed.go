package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// feed adds the "feed" commands.
func (r Router) feed(root *cobra.Command) {
	feed := group(root, "operations", "feed", "The feeds: the catalog API, directories, Assembly64 …", "feeds")
	add(feed, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the feeds, the catalog API first", Args: ExactArgs(0),
		Example: "  warp-engine feed list\n  warp-engine feed list -o json",
		Run:     r.FeedController.Index,
	})
	add(feed, "", Route{
		Use: "check KEY", Short: "Show what a feed offers and which softwares it leaves out", Args: ExactArgs(1, "KEY"),
		Example: "  warp-engine feed check assembly64\n  warp-engine feed check c64 --strict",
		Flags: func(flags *pflag.FlagSet) {
			flags.Bool("strict", false, "fail (exit code 1) when the feed leaves a software out")
		},
		Run: r.FeedController.Check,
	})
}
