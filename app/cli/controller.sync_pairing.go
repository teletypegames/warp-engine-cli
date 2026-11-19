package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SyncPairingControllerInterface is the "pairing" commands: the inbound
// peers the caller paired with the engine itself.
type SyncPairingControllerInterface interface {
	// Index lists the pairings.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one pairing.
	Show(cmd *cobra.Command, args []string) error
	// Create pairs and prints the new token.
	Create(cmd *cobra.Command, args []string) error
	// Update renames a pairing.
	Update(cmd *cobra.Command, args []string) error
	// Destroy unpairs.
	Destroy(cmd *cobra.Command, args []string) error
	// Token issues a new token for a pairing.
	Token(cmd *cobra.Command, args []string) error
}

// SyncPairingController is the [SyncPairingControllerInterface] over a
// [domain.SyncPairingServiceInterface].
type SyncPairingController struct {
	// SyncPairingService manages the pairings.
	SyncPairingService domain.SyncPairingServiceInterface
}

// pairingHint tells what to do with a pairing's new token.
const pairingHint = "Set it as the secret of the other engine's sync peer for this one."

// Index prints the pairings, one page or, with --all, every page.
func (c SyncPairingController) Index(cmd *cobra.Command, args []string) error {
	serializer := SyncPairingSerializer{}
	page, err := c.SyncPairingService.All(domain.Filter{}, FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the pairing ID.
func (c SyncPairingController) Show(cmd *cobra.Command, args []string) error {
	sync_pairing_id, err := argID("pairing", args[0])
	if err != nil {
		return err
	}
	pairing, err := c.SyncPairingService.Show(sync_pairing_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, pairing, SyncPairingSerializer{}.ModelToRecord(*pairing))
}

// Create pairs with --label and prints the pairing and its new token, the
// only time the whole token is shown.
func (c SyncPairingController) Create(cmd *cobra.Command, args []string) error {
	serializer := SyncPairingSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	pairing, err := c.SyncPairingService.Create(attributes)
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, pairing, serializer.ModelToRecord(*pairing)); err != nil {
		return err
	}
	printPlainToken(cmd, pairing.PlainToken, pairingHint)
	return nil
}

// Update renames the pairing ID with --label and prints it.
func (c SyncPairingController) Update(cmd *cobra.Command, args []string) error {
	serializer := SyncPairingSerializer{}
	sync_pairing_id, err := argID("pairing", args[0])
	if err != nil {
		return err
	}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	pairing, err := c.SyncPairingService.Update(sync_pairing_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, pairing, serializer.ModelToRecord(*pairing))
}

// Destroy unpairs the pairing ID after asking (not with --yes).
func (c SyncPairingController) Destroy(cmd *cobra.Command, args []string) error {
	sync_pairing_id, err := argID("pairing", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Unpair "+args[0]+" and revoke its token?"); err != nil {
		return err
	}
	if err = c.SyncPairingService.Destroy(sync_pairing_id); err != nil {
		return err
	}
	PrintDone(cmd, "Pairing %d deleted.", sync_pairing_id)
	return nil
}

// Token issues a new token for the pairing ID after asking (the old one
// stops working) and prints it, the only time it is shown.
func (c SyncPairingController) Token(cmd *cobra.Command, args []string) error {
	sync_pairing_id, err := argID("pairing", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Issue a new token for the pairing "+args[0]+"? Its current one stops working."); err != nil {
		return err
	}
	pairing, err := c.SyncPairingService.Token(sync_pairing_id)
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, pairing, SyncPairingSerializer{}.ModelToRecord(*pairing)); err != nil {
		return err
	}
	printPlainToken(cmd, pairing.PlainToken, pairingHint)
	return nil
}
