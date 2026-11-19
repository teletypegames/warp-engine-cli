package cli

import (
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SyncEventSerializer turns the flags of the sync-event commands into the
// API's terms and the sync events and states into tables and records.
type SyncEventSerializer struct{}

// SyncEventSorts are the orders "sync-event list --sort" takes.
var SyncEventSorts = []string{"createdAt"}

// Filters are the filter flags of "sync-event list". --software is handled
// by [SyncEventSerializer.FlagsToFilter]: an ID filters on softwareId, a
// name on software.
func (s SyncEventSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "peer", Key: "syncPeerId", Usage: "sync peer ID is"},
		{Flag: "direction", Key: "direction", Usage: "direction is (" + strings.Join(domain.SyncDirections, ", ") + ")"},
		{Flag: "status", Key: "status", Usage: "status is (" + strings.Join(domain.SyncEventStatuses, ", ") + ")"},
		{Flag: "software", Key: "software", Usage: "software is (an ID or a name)"},
	}
}

// FlagsToFilter reads the filter flags of "sync-event list"; a numeric
// --software becomes softwareId.
func (s SyncEventSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	filter := FlagsToFilter(flags, s.Filters())
	if software, ok := filter["software"]; ok {
		if _, err := strconv.Atoi(software); err == nil {
			delete(filter, "software")
			filter["softwareId"] = software
		}
	}
	return filter
}

// ModelToTable renders sync events as a table: ID, time, direction,
// status, peer, software, version, step, message.
func (s SyncEventSerializer) ModelToTable(events []domain.SyncEventModel) Table {
	table := Table{Headers: []string{"ID", "When", "Direction", "Status", "Peer", "Software", "Version", "Step", "Message"}}
	for _, event := range events {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(event.ID), Text(event.CreatedAt), event.Direction, event.Status, Number(event.SyncPeerID),
			Text(event.SoftwareName), Text(event.Version), Text(event.Step), Truncate(Text(event.Message), 50),
		})
	}
	return table
}

// ModelToRecord renders one sync event with every field.
func (s SyncEventSerializer) ModelToRecord(event domain.SyncEventModel) Record {
	size := "-"
	if event.Bytes != nil {
		size = Bytes(*event.Bytes)
	}
	duration := "-"
	if event.DurationMs != nil {
		duration = strconv.Itoa(*event.DurationMs) + " ms"
	}
	return Record{
		{Name: "ID", Value: strconv.Itoa(event.ID)},
		{Name: "Direction", Value: event.Direction},
		{Name: "Status", Value: event.Status},
		{Name: "Peer", Value: Number(event.SyncPeerID)},
		{Name: "Software", Value: Text(event.SoftwareName)},
		{Name: "Software ID", Value: Number(event.SoftwareID)},
		{Name: "Platform", Value: Text(event.Platform)},
		{Name: "Version", Value: Text(event.Version)},
		{Name: "Release", Value: Number(event.ReleaseID)},
		{Name: "Step", Value: Text(event.Step)},
		{Name: "HTTP status", Value: Number(event.HTTPStatus)},
		{Name: "Message", Value: Text(event.Message)},
		{Name: "Files", Value: Number(event.FilesCount)},
		{Name: "Size", Value: size},
		{Name: "Duration", Value: duration},
		{Name: "Attempt", Value: Number(event.Attempt)},
		{Name: "Token", Value: Text(event.TokenPrefix)},
		{Name: "Created", Value: Text(event.CreatedAt)},
		{Name: "Updated", Value: Text(event.UpdatedAt)},
	}
}

// StatusToTable renders the sync states of a software: peer ID, peer,
// direction, then the latest event's status, version and time, and the
// reason the newest release would not go.
func (s SyncEventSerializer) StatusToTable(states []domain.SyncStatusModel) Table {
	table := Table{Headers: []string{"Peer ID", "Peer", "Direction", "Status", "Version", "When", "Reason"}}
	for _, state := range states {
		status, version, when := "-", "-", "-"
		if state.Event != nil {
			status, version, when = state.Event.Status, Text(state.Event.Version), Text(state.Event.CreatedAt)
		}
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(state.PeerID), state.Peer, state.Direction, status, version, when, Text(state.Reason),
		})
	}
	return table
}
