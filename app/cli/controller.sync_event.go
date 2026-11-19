package cli

import (
	"fmt"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SyncEventControllerInterface is the "sync-event" commands: the pushes
// between this engine and its peers, and where each software stands.
type SyncEventControllerInterface interface {
	// Index lists the sync events.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one sync event.
	Show(cmd *cobra.Command, args []string) error
	// Software lists the sync events of one software.
	Software(cmd *cobra.Command, args []string) error
	// Status shows where one software stands with every peer.
	Status(cmd *cobra.Command, args []string) error
}

// SyncEventController is the [SyncEventControllerInterface] over a
// [domain.SyncEventServiceInterface].
type SyncEventController struct {
	// SyncEventService reads the sync events.
	SyncEventService domain.SyncEventServiceInterface
}

// Index prints the sync events matching --peer, --software, --direction
// and --status, newest first, one page or, with --all, every page.
func (c SyncEventController) Index(cmd *cobra.Command, args []string) error {
	serializer := SyncEventSerializer{}
	page, err := c.SyncEventService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the sync event ID.
func (c SyncEventController) Show(cmd *cobra.Command, args []string) error {
	sync_event_id, err := argID("sync event", args[0])
	if err != nil {
		return err
	}
	event, err := c.SyncEventService.Show(sync_event_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, event, SyncEventSerializer{}.ModelToRecord(*event))
}

// Software prints the sync events of SOFTWARE (an ID or a name), in the
// --direction given or both.
func (c SyncEventController) Software(cmd *cobra.Command, args []string) error {
	direction, _ := cmd.Flags().GetString("direction")
	if direction != "" && direction != "outbound" && direction != "inbound" {
		return &UsageError{Err: fmt.Errorf("--direction takes one of %s, not %q", strings.Join(domain.SyncDirections, ", "), direction)}
	}
	events, err := c.SyncEventService.Software(args[0], direction)
	if err != nil {
		return err
	}
	return PrintList(cmd, events, SyncEventSerializer{}.ModelToTable(events))
}

// Status prints where SOFTWARE (an ID or a name) stands with every peer.
func (c SyncEventController) Status(cmd *cobra.Command, args []string) error {
	states, err := c.SyncEventService.Status(args[0])
	if err != nil {
		return err
	}
	return PrintList(cmd, states, SyncEventSerializer{}.StatusToTable(states))
}
