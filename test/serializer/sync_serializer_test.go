package serializer_test

import (
	"testing"

	"git.teletypegames.org/engines/warp-engine-cli/app/cli"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// peerFlags declares the peer create/update flags and parses args.
func peerFlags(t *testing.T, args ...string) *pflag.FlagSet {
	t.Helper()
	serializer := cli.SyncPeerSerializer{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFields(flags, serializer.Fields())
	serializer.DeclareOptions(flags)
	assert.NoError(t, flags.Parse(args))
	return flags
}

func TestSyncPeerSerializer_RequestToAttributes_OptionsOverData(t *testing.T) {
	flags := peerFlags(t, "--data", `{"transportOptions":{"host":"old","port":"22"}}`, "--option", "host=new", "--option", "path=/srv/a=b")

	attributes, err := cli.SyncPeerSerializer{}.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"transportOptions": map[string]any{"host": "new", "port": "22", "path": "/srv/a=b"}}, attributes)
}

func TestSyncPeerSerializer_RequestToAttributes_OptionWithoutEqualsIsUsageError(t *testing.T) {
	flags := peerFlags(t, "--option", "host")

	_, err := cli.SyncPeerSerializer{}.RequestToAttributes(flags)

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSyncPeerSerializer_RequestToAttributes_NoOptionsSendsNone(t *testing.T) {
	flags := peerFlags(t, "--label", "games", "--inbound")

	attributes, err := cli.SyncPeerSerializer{}.RequestToAttributes(flags)

	assert.NoError(t, err)
	assert.Equal(t, domain.Attributes{"label": "games", "inboundEnabled": true}, attributes)
}

func TestSyncPeerSerializer_RequestToAttributes_UnknownChannelIsUsageError(t *testing.T) {
	flags := peerFlags(t, "--channels", "nightly")

	_, err := cli.SyncPeerSerializer{}.RequestToAttributes(flags)

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestSyncPeerSerializer_ModelToTable_MapsFields(t *testing.T) {
	table := cli.SyncPeerSerializer{}.ModelToTable([]domain.SyncPeerModel{{
		ID: 3, Label: "games", Transport: "http", TransportOptions: map[string]any{"host": "https://games.example.org"},
		OutboundEnabled: true, AllowedChannels: []string{"stable"},
	}})

	assert.Equal(t, []string{"3", "games", "http", "https://games.example.org", "yes", "no", "stable", "-", "-"}, table.Rows[0])
}

func TestSyncPeerSerializer_ModelToTable_Empty(t *testing.T) {
	table := cli.SyncPeerSerializer{}.ModelToTable(nil)

	assert.Empty(t, table.Rows)
	assert.NotEmpty(t, table.Headers)
}

func TestSyncPeerSerializer_ModelToRecord_IDFirstAndOptionsSorted(t *testing.T) {
	record := cli.SyncPeerSerializer{}.ModelToRecord(domain.SyncPeerModel{ID: 3, TransportOptions: map[string]any{"user": "warp", "host": "h", "port": 22}})

	assert.Equal(t, cli.Field{Name: "ID", Value: "3"}, record[0])
	assert.Equal(t, cli.Field{Name: "Option host", Value: "h"}, record[3])
	assert.Equal(t, cli.Field{Name: "Option port", Value: "22"}, record[4])
	assert.Equal(t, cli.Field{Name: "Option user", Value: "warp"}, record[5])
}

func TestSyncPeerSerializer_PushToRecord_ListsReleasesAndSkipped(t *testing.T) {
	record := cli.SyncPeerSerializer{}.PushToRecord(domain.SyncPushModel{
		Queued: 1, Peers: []int{3, 5}, Releases: []domain.SyncPushReleaseModel{{Software: "pong", Version: "1.0.0"}}, Skipped: []string{"ra 2.0: beta"},
	})

	assert.Contains(t, record, cli.Field{Name: "Peers", Value: "3, 5"})
	assert.Contains(t, record, cli.Field{Name: "Release", Value: "pong 1.0.0"})
	assert.Contains(t, record, cli.Field{Name: "Skipped", Value: "ra 2.0: beta"})
}

func TestSyncPairingSerializer_ModelToTable_FilledAndEmpty(t *testing.T) {
	prefix := "9f8e"
	table := cli.SyncPairingSerializer{}.ModelToTable([]domain.SyncPairingModel{{ID: 2, Label: "laptop", InboundTokenPrefix: &prefix}})
	assert.Equal(t, []string{"2", "laptop", "9f8e", "-", "-"}, table.Rows[0])

	empty := cli.SyncPairingSerializer{}.ModelToTable(nil)
	assert.Empty(t, empty.Rows)
	assert.NotEmpty(t, empty.Headers)
}

func TestSyncPairingSerializer_ModelToRecord_HasNoPlainToken(t *testing.T) {
	plain := "secret"
	record := cli.SyncPairingSerializer{}.ModelToRecord(domain.SyncPairingModel{ID: 2, PlainToken: &plain})

	assert.Equal(t, "ID", record[0].Name)
	for _, field := range record {
		assert.NotEqual(t, "secret", field.Value)
	}
}

func TestSyncEventSerializer_FlagsToFilter_NumericSoftwareIsAnID(t *testing.T) {
	serializer := cli.SyncEventSerializer{}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFilters(flags, serializer.Filters())
	assert.NoError(t, flags.Parse([]string{"--software", "12", "--direction", "inbound"}))

	assert.Equal(t, domain.Filter{"softwareId": "12", "direction": "inbound"}, serializer.FlagsToFilter(flags))
}

func TestSyncEventSerializer_ModelToTable_FilledAndEmpty(t *testing.T) {
	name, version := "pong", "1.0.0"
	peer := 3
	table := cli.SyncEventSerializer{}.ModelToTable([]domain.SyncEventModel{{ID: 9, Direction: "outbound", Status: "succeeded", SyncPeerID: &peer, SoftwareName: &name, Version: &version}})
	assert.Equal(t, []string{"9", "-", "outbound", "succeeded", "3", "pong", "1.0.0", "-", "-"}, table.Rows[0])

	empty := cli.SyncEventSerializer{}.ModelToTable(nil)
	assert.Empty(t, empty.Rows)
	assert.NotEmpty(t, empty.Headers)
}

func TestSyncEventSerializer_ModelToRecord_RendersSizeAndDuration(t *testing.T) {
	size := int64(2048)
	duration := 1500
	record := cli.SyncEventSerializer{}.ModelToRecord(domain.SyncEventModel{ID: 9, Bytes: &size, DurationMs: &duration})

	assert.Equal(t, cli.Field{Name: "ID", Value: "9"}, record[0])
	assert.Contains(t, record, cli.Field{Name: "Size", Value: "2.0 KiB"})
	assert.Contains(t, record, cli.Field{Name: "Duration", Value: "1500 ms"})
}

func TestSyncEventSerializer_StatusToTable_WithAndWithoutEvent(t *testing.T) {
	version, reason := "1.2.0", "not allowed"
	table := cli.SyncEventSerializer{}.StatusToTable([]domain.SyncStatusModel{
		{PeerID: 3, Peer: "games", Direction: "outbound", Event: &domain.SyncEventModel{Status: "succeeded", Version: &version}},
		{PeerID: 4, Peer: "mirror", Direction: "outbound", Reason: &reason},
	})

	assert.Equal(t, []string{"3", "games", "outbound", "succeeded", "1.2.0", "-", "-"}, table.Rows[0])
	assert.Equal(t, []string{"4", "mirror", "outbound", "-", "-", "-", "not allowed"}, table.Rows[1])
	assert.Empty(t, cli.SyncEventSerializer{}.StatusToTable(nil).Rows)
}

func TestApplicationTokenSerializer_CreateToAttributes_OwnerIDIsANumber(t *testing.T) {
	serializer := cli.ApplicationTokenSerializer{}
	flags := parse(t, serializer.CreateFields(), "--name", "ci", "--scopes", "admin:read", "--owner-id", "x1")

	_, err := serializer.CreateToAttributes(flags)

	assert.Equal(t, cli.ExitUsage, cli.ExitCode(err))
}

func TestApplicationTokenSerializer_RequestToAttributes_HasNoOwnerID(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	cli.DeclareFields(flags, cli.ApplicationTokenSerializer{}.Fields())

	assert.Error(t, flags.Parse([]string{"--owner-id", "1"}))
}

func TestApplicationTokenSerializer_ModelToTable_FilledAndEmpty(t *testing.T) {
	table := cli.ApplicationTokenSerializer{}.ModelToTable([]domain.ApplicationTokenModel{{ID: 2, Name: "ci", TokenPrefix: "304f", Scopes: []string{"update", "upload"}}})
	assert.Equal(t, []string{"2", "ci", "304f", "update, upload", "no", "-", "-"}, table.Rows[0])

	empty := cli.ApplicationTokenSerializer{}.ModelToTable(nil)
	assert.Empty(t, empty.Rows)
	assert.NotEmpty(t, empty.Headers)
}

func TestApplicationTokenSerializer_ModelToRecord_SyncCountsAndNoPlainToken(t *testing.T) {
	plain := "whole"
	record := cli.ApplicationTokenSerializer{}.ModelToRecord(domain.ApplicationTokenModel{ID: 5, PlainToken: &plain, Sync: &domain.ApplicationTokenSyncModel{Synced: 2, Failed: 1}})

	assert.Equal(t, cli.Field{Name: "ID", Value: "5"}, record[0])
	assert.Contains(t, record, cli.Field{Name: "CI secrets", Value: "2 updated, 1 failed"})
	for _, field := range record {
		assert.NotEqual(t, "whole", field.Value)
	}
}
