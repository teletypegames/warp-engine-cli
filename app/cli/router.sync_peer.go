package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// syncPeer adds the "peer" commands.
func (r Router) syncPeer(root *cobra.Command) {
	serializer := SyncPeerSerializer{}
	peer := group(root, "sync", "peer", "The other engines this one syncs releases with", "peers")
	add(peer, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the sync peers", Args: ExactArgs(0),
		Example: "  warp-engine peer list --outbound\n  warp-engine peer list --transport rsync -o json",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, SyncPeerSorts)
		},
		Run: r.SyncPeerController.Index,
	})
	add(peer, "", Route{
		Use: "show ID", Short: "Show a sync peer", Args: ExactArgs(1, "ID"),
		Run: r.SyncPeerController.Show,
	})
	add(peer, "", Route{
		Use: "create", Short: "Add a sync peer", Args: ExactArgs(0),
		Example: "  warp-engine peer create --label games --option host=https://games.example.org --secret we_xxx --outbound\n" +
			"  warp-engine peer create --label mirror --transport rsync --option host=mirror.example.org --option user=warp --option path=/srv/warp",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFields(flags, serializer.Fields())
			serializer.DeclareOptions(flags)
		},
		Run: r.SyncPeerController.Create,
	})
	add(peer, "", Route{
		Use: "update ID", Short: "Change the fields given as flags", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine peer update 3 --outbound=false\n  warp-engine peer update 3 --channels stable,beta",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFields(flags, serializer.Fields())
			serializer.DeclareOptions(flags)
		},
		Run: r.SyncPeerController.Update,
	})
	add(peer, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a sync peer and revoke its inbound token", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SyncPeerController.Destroy,
	})
	add(peer, "", Route{
		Use: "check ID", Short: "Ask whether the peer accepts this engine", Args: ExactArgs(1, "ID"),
		Long: "Ask whether the peer accepts this engine: it is reachable and takes the secret.\nExits 1 when it does not.",
		Run:  r.SyncPeerController.Check,
	})
	add(peer, "", Route{
		Use: "inbound-token ID", Short: "Issue a new token the peer pushes here with (shown once)", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SyncPeerController.InboundToken,
	})
	add(peer, "", Route{
		Use: "push ID", Short: "Push the releases the peer may be missing", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine peer push 3 --dry-run\n  warp-engine peer push 3 --software blessing-of-ra --channels stable --yes",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFields(flags, serializer.PushFields())
			DeclareYes(flags)
		},
		Run: r.SyncPeerController.Push,
	})
	add(peer, "", Route{
		Use: "run", Short: "Push to every outbound peer, or to the --peer ones", Args: ExactArgs(0),
		Example: "  warp-engine peer run --dry-run\n  warp-engine peer run --peer 3 --peer 5 --all --yes",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFields(flags, serializer.PushFields())
			serializer.DeclarePeers(flags)
			DeclareYes(flags)
		},
		Run: r.SyncPeerController.Run,
	})
}
