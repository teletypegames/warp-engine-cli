package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// syncPairing adds the "pairing" commands.
func (r Router) syncPairing(root *cobra.Command) {
	serializer := SyncPairingSerializer{}
	pairing := group(root, "sync", "pairing", "Your own pairings: other engines you let push here", "pairings")
	add(pairing, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List your pairings", Args: ExactArgs(0),
		Flags: func(flags *pflag.FlagSet) { DeclareListFlags(flags, nil) },
		Run:   r.SyncPairingController.Index,
	})
	add(pairing, "", Route{
		Use: "show ID", Short: "Show a pairing", Args: ExactArgs(1, "ID"),
		Run: r.SyncPairingController.Show,
	})
	add(pairing, "", Route{
		Use: "create", Short: "Pair: add an inbound peer and issue its token (shown once)", Args: ExactArgs(0),
		Example: "  warp-engine pairing create --label \"laptop local\"",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SyncPairingController.Create,
	})
	add(pairing, "", Route{
		Use: "update ID", Short: "Rename a pairing", Args: ExactArgs(1, "ID"),
		Flags: func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:   r.SyncPairingController.Update,
	})
	add(pairing, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Unpair and revoke the token", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SyncPairingController.Destroy,
	})
	add(pairing, "", Route{
		Use: "token ID", Short: "Issue a new token for a pairing (shown once)", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SyncPairingController.Token,
	})
}
