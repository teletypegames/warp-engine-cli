package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// syncEvent adds the "sync-event" commands.
func (r Router) syncEvent(root *cobra.Command) {
	serializer := SyncEventSerializer{}
	event := group(root, "sync", "sync-event", "What was pushed to and taken from the sync peers", "sync-events")
	add(event, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the sync events, newest first", Args: ExactArgs(0),
		Example: "  warp-engine sync-event list --status failed\n  warp-engine sync-event list --peer 3 --software blessing-of-ra",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, SyncEventSorts)
		},
		Run: r.SyncEventController.Index,
	})
	add(event, "", Route{
		Use: "show ID", Short: "Show a sync event", Args: ExactArgs(1, "ID"),
		Run: r.SyncEventController.Show,
	})
	add(event, "", Route{
		Use: "software SOFTWARE", Short: "List the sync events of a software (by ID or name)", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine sync-event software blessing-of-ra --direction outbound",
		Flags: func(flags *pflag.FlagSet) {
			flags.String("direction", "", "outbound or inbound (default both)")
		},
		Run: r.SyncEventController.Software,
	})
	add(event, "", Route{
		Use: "status SOFTWARE", Short: "Show where a software stands with every peer", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine sync-event status blessing-of-ra",
		Run:     r.SyncEventController.Status,
	})
}
