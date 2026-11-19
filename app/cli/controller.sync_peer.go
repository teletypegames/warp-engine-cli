package cli

import (
	"errors"
	"fmt"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SyncPeerControllerInterface is the "peer" commands: the other engines
// this one syncs with.
type SyncPeerControllerInterface interface {
	// Index lists the sync peers.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one sync peer.
	Show(cmd *cobra.Command, args []string) error
	// Create adds a sync peer.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a sync peer.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a sync peer.
	Destroy(cmd *cobra.Command, args []string) error
	// Check asks whether a sync peer accepts this engine.
	Check(cmd *cobra.Command, args []string) error
	// InboundToken issues a new inbound token for a sync peer.
	InboundToken(cmd *cobra.Command, args []string) error
	// Push queues the releases a sync peer may be missing.
	Push(cmd *cobra.Command, args []string) error
	// Run queues the pushes to several sync peers.
	Run(cmd *cobra.Command, args []string) error
}

// SyncPeerController is the [SyncPeerControllerInterface] over a
// [domain.SyncPeerServiceInterface].
type SyncPeerController struct {
	// SyncPeerService manages the sync peers.
	SyncPeerService domain.SyncPeerServiceInterface
}

// ErrCheckFailed is returned by "peer check" when the peer did not accept
// this engine; the answer is printed before.
var ErrCheckFailed = errors.New("the peer did not accept the check")

// Index prints the sync peers matching --outbound, --inbound and
// --transport, one page or, with --all, every page.
func (c SyncPeerController) Index(cmd *cobra.Command, args []string) error {
	serializer := SyncPeerSerializer{}
	page, err := c.SyncPeerService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the sync peer ID.
func (c SyncPeerController) Show(cmd *cobra.Command, args []string) error {
	sync_peer_id, err := argID("sync peer", args[0])
	if err != nil {
		return err
	}
	peer, err := c.SyncPeerService.Show(sync_peer_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, peer, SyncPeerSerializer{}.ModelToRecord(*peer))
}

// Create adds a sync peer from the field flags and --option (--label is
// needed) and prints it.
func (c SyncPeerController) Create(cmd *cobra.Command, args []string) error {
	serializer := SyncPeerSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	peer, err := c.SyncPeerService.Create(attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, peer, serializer.ModelToRecord(*peer))
}

// Update changes the fields given as flags on the sync peer ID and prints it.
func (c SyncPeerController) Update(cmd *cobra.Command, args []string) error {
	serializer := SyncPeerSerializer{}
	sync_peer_id, err := argID("sync peer", args[0])
	if err != nil {
		return err
	}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	peer, err := c.SyncPeerService.Update(sync_peer_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, peer, serializer.ModelToRecord(*peer))
}

// Destroy deletes the sync peer ID after asking (not with --yes).
func (c SyncPeerController) Destroy(cmd *cobra.Command, args []string) error {
	sync_peer_id, err := argID("sync peer", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Delete the sync peer "+args[0]+" and revoke its inbound token?"); err != nil {
		return err
	}
	if err = c.SyncPeerService.Destroy(sync_peer_id); err != nil {
		return err
	}
	PrintDone(cmd, "Sync peer %d deleted.", sync_peer_id)
	return nil
}

// Check asks the sync peer ID whether it accepts this engine and prints
// the answer; a refusal is [ErrCheckFailed].
func (c SyncPeerController) Check(cmd *cobra.Command, args []string) error {
	sync_peer_id, err := argID("sync peer", args[0])
	if err != nil {
		return err
	}
	check, err := c.SyncPeerService.Check(sync_peer_id)
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, check, SyncPeerSerializer{}.CheckToRecord(*check)); err != nil {
		return err
	}
	if !check.OK {
		return ErrCheckFailed
	}
	return nil
}

// InboundToken issues a new inbound token for the sync peer ID after
// asking (the old one stops working) and prints it: the only time the
// whole token is shown.
func (c SyncPeerController) InboundToken(cmd *cobra.Command, args []string) error {
	sync_peer_id, err := argID("sync peer", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Issue a new inbound token for the sync peer "+args[0]+"? Its current one stops working."); err != nil {
		return err
	}
	token, err := c.SyncPeerService.InboundToken(sync_peer_id)
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, token, SyncPeerSerializer{}.InboundTokenToRecord(*token)); err != nil {
		return err
	}
	printPlainToken(cmd, &token.PlainToken, "Give it to the other engine as the secret of its sync peer for this one.")
	return nil
}

// Push queues the releases the sync peer ID may be missing, chosen by the
// push flags, after asking (not with --yes or --dry-run), and prints what goes.
func (c SyncPeerController) Push(cmd *cobra.Command, args []string) error {
	serializer := SyncPeerSerializer{}
	sync_peer_id, err := argID("sync peer", args[0])
	if err != nil {
		return err
	}
	attributes, err := serializer.PushToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	if dry_run, _ := attributes["dryRun"].(bool); !dry_run {
		if err = confirmOrCancel(cmd, "Push releases to the sync peer "+args[0]+"?"); err != nil {
			return err
		}
	}
	push, err := c.SyncPeerService.Push(sync_peer_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, push, serializer.PushToRecord(*push))
}

// Run queues the pushes to every outbound peer, or to the --peer ones,
// after asking (not with --yes or --dry-run), and prints what goes.
func (c SyncPeerController) Run(cmd *cobra.Command, args []string) error {
	serializer := SyncPeerSerializer{}
	attributes, err := serializer.RunToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	if dry_run, _ := attributes["dryRun"].(bool); !dry_run {
		question := "Push releases to every outbound sync peer?"
		if attributes.Has("peerIds") {
			question = fmt.Sprintf("Push releases to the sync peers %v?", attributes["peerIds"])
		}
		if err = confirmOrCancel(cmd, question); err != nil {
			return err
		}
	}
	push, err := c.SyncPeerService.Run(attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, push, serializer.PushToRecord(*push))
}

// argID parses the ID argument text of a kind of record ("sync peer"); a
// malformed one is a [*UsageError].
func argID(kind string, text string) (int, error) {
	id, err := domain.ParseID(kind, text)
	if err != nil {
		return 0, &UsageError{Err: err}
	}
	return id, nil
}

// confirmOrCancel asks question with [Confirm]; a no is [ErrCancelled].
func confirmOrCancel(cmd *cobra.Command, question string) error {
	ok, err := Confirm(cmd, question)
	if err == nil && !ok {
		err = ErrCancelled
	}
	return err
}

// printPlainToken prints a token the server shows only once, set apart
// below the record, with hint, for [OutputTable]; the JSON output carries
// it as plainToken.
func printPlainToken(cmd *cobra.Command, plain *string, hint string) {
	if outputOf(cmd) != OutputTable || plain == nil || *plain == "" {
		return
	}
	writer := cmd.OutOrStdout()
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "New token, shown only this once — store it now:")
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "    "+*plain)
	fmt.Fprintln(writer)
	if hint != "" {
		fmt.Fprintln(writer, hint)
	}
}
