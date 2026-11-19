package cli

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SyncPeerSerializer turns the flags of the peer commands into the API's
// terms and the sync peers, checks and pushes into tables and records.
type SyncPeerSerializer struct{}

// SyncPeerSorts are the orders "peer list --sort" takes.
var SyncPeerSorts = []string{"label", "lastSyncedAt"}

// OptionFlag is the repeatable flag of "peer create" and "peer update"
// that sets one transport option: --option host=https://games.example.org.
const OptionFlag = "option"

// PeerFlag is the repeatable flag of "peer run" naming a peer to push to.
const PeerFlag = "peer"

// Filters are the filter flags of "peer list".
func (s SyncPeerSerializer) Filters() []FilterFlag {
	return []FilterFlag{
		{Flag: "outbound", Key: "outboundEnabled", Usage: "pushed to or not", Bool: true},
		{Flag: "inbound", Key: "inboundEnabled", Usage: "accepted from or not", Bool: true},
		{Flag: "transport", Key: "transport", Usage: "transport is (" + strings.Join(domain.SyncTransports, ", ") + ")"},
	}
}

// Fields are the field flags of "peer create" and "peer update"; the
// transport options come with --option ([OptionFlag]).
func (s SyncPeerSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "label", Key: "label", Usage: "display name"},
		{Flag: "transport", Key: "transport", Usage: "how it is reached (default http)", Enum: domain.SyncTransports},
		{Flag: "secret", Key: "secret", Usage: "the transport's secret (http: the token the other engine issued); write-only"},
		{Flag: "outbound", Key: "outboundEnabled", Kind: KindBool, Usage: "push our releases to it"},
		{Flag: "inbound", Key: "inboundEnabled", Kind: KindBool, Usage: "accept its pushes"},
		{Flag: "channels", Key: "allowedChannels", Kind: KindStrings, Usage: "channels pushed to it; none: every allowed one", Enum: domain.Channels},
	}
}

// PushFields are the field flags of "peer push" and "peer run": which
// releases go, and --dry-run.
func (s SyncPeerSerializer) PushFields() []FieldFlag {
	return []FieldFlag{
		{Flag: "all", Key: "all", Kind: KindBool, Usage: "every allowed release, oldest first (a full sync); default: the newest of each software"},
		{Flag: "filtered", Key: "filtered", Kind: KindBool, Usage: "apply the configured filters (default true; --filtered=false sends everything)"},
		{Flag: "channels", Key: "channels", Kind: KindStrings, Usage: "channels (default: the configured ones)", Enum: domain.Channels},
		{Flag: "softwares", Key: "softwares", Kind: KindStrings, Usage: "software names, * as a wildcard"},
		{Flag: "software", Key: "software", Usage: "one software, by name"},
		{Flag: "platforms", Key: "platforms", Kind: KindStrings, Usage: "platform names"},
		{Flag: "categories", Key: "categories", Kind: KindStrings, Usage: "category keys"},
		{Flag: "statuses", Key: "statuses", Kind: KindStrings, Usage: "status keys"},
		{Flag: "skip-existing", Key: "skipExisting", Kind: KindBool, Usage: "leave out what the peer already has (default true)"},
		{Flag: "dry-run", Key: "dryRun", Kind: KindBool, Usage: "only show what would go"},
	}
}

// DeclareOptions adds the repeatable --option flag.
func (s SyncPeerSerializer) DeclareOptions(flags *pflag.FlagSet) {
	flags.StringArray(OptionFlag, nil, "a transport option as key=value, repeatable (http: host; rsync: host, user, path, port, host_key); replaces the stored options")
}

// DeclarePeers adds the repeatable --peer flag of "peer run".
func (s SyncPeerSerializer) DeclarePeers(flags *pflag.FlagSet) {
	flags.IntSlice(PeerFlag, nil, "a peer ID to push to, repeatable (default: every outbound peer)")
}

// FlagsToFilter reads the filter flags of "peer list".
func (s SyncPeerSerializer) FlagsToFilter(flags *pflag.FlagSet) domain.Filter {
	return FlagsToFilter(flags, s.Filters())
}

// RequestToAttributes reads the field flags of create and update, and the
// --option flags into transportOptions (over the --data object's
// transportOptions). An option without "=" is a [*UsageError].
func (s SyncPeerSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	attributes, err := FieldsToAttributes(flags, s.Fields())
	if err != nil {
		return nil, err
	}
	if !flags.Changed(OptionFlag) {
		return attributes, nil
	}
	pairs, err := flags.GetStringArray(OptionFlag)
	if err != nil {
		return nil, err
	}
	options := map[string]any{}
	if given, ok := attributes["transportOptions"].(map[string]any); ok {
		for key, value := range given {
			options[key] = value
		}
	}
	for _, pair := range pairs {
		key, value, found := strings.Cut(pair, "=")
		if !found || key == "" {
			return nil, &UsageError{Err: fmt.Errorf("--%s takes key=value, not %q", OptionFlag, pair)}
		}
		options[key] = value
	}
	attributes["transportOptions"] = options
	return attributes, nil
}

// PushToAttributes reads the flags of "peer push": the push criteria.
func (s SyncPeerSerializer) PushToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.PushFields())
}

// RunToAttributes reads the flags of "peer run": the push criteria and the
// --peer IDs as peerIds.
func (s SyncPeerSerializer) RunToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	attributes, err := s.PushToAttributes(flags)
	if err != nil {
		return nil, err
	}
	if flags.Changed(PeerFlag) {
		peer_ids, err := flags.GetIntSlice(PeerFlag)
		if err != nil {
			return nil, err
		}
		attributes["peerIds"] = peer_ids
	}
	return attributes, nil
}

// ModelToTable renders sync peers as a table: ID, label, transport, host,
// outbound, inbound, channels, last synced, last error.
func (s SyncPeerSerializer) ModelToTable(peers []domain.SyncPeerModel) Table {
	table := Table{Headers: []string{"ID", "Label", "Transport", "Host", "Outbound", "Inbound", "Channels", "Last synced", "Last error"}}
	for _, peer := range peers {
		host := "-"
		if value, ok := peer.TransportOptions["host"]; ok {
			host = fmt.Sprint(value)
		}
		last_error := Text(peer.LastError)
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(peer.ID), peer.Label, peer.Transport, host, YesNo(peer.OutboundEnabled), YesNo(peer.InboundEnabled),
			List(peer.AllowedChannels), Text(peer.LastSyncedAt), Truncate(last_error, 40),
		})
	}
	return table
}

// ModelToRecord renders one sync peer with every field; each transport
// option is a line "Option <key>".
func (s SyncPeerSerializer) ModelToRecord(peer domain.SyncPeerModel) Record {
	record := Record{
		{Name: "ID", Value: strconv.Itoa(peer.ID)},
		{Name: "Label", Value: peer.Label},
		{Name: "Transport", Value: peer.Transport},
	}
	keys := make([]string, 0, len(peer.TransportOptions))
	for key := range peer.TransportOptions {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := peer.TransportOptions[key]
		text, ok := value.(string)
		if !ok {
			text = JSONText(value)
		}
		record = append(record, Field{Name: "Option " + key, Value: text})
	}
	created_by := "-"
	if peer.CreatedBy != nil {
		created_by = peer.CreatedBy.Type + " #" + strconv.Itoa(peer.CreatedBy.ID)
	}
	return append(record,
		Field{Name: "Secret", Value: Text(peer.SecretPrefix)},
		Field{Name: "Outbound", Value: YesNo(peer.OutboundEnabled)},
		Field{Name: "Inbound", Value: YesNo(peer.InboundEnabled)},
		Field{Name: "Channels", Value: List(peer.AllowedChannels)},
		Field{Name: "Inbound token", Value: Text(peer.InboundTokenPrefix)},
		Field{Name: "Paired by", Value: created_by},
		Field{Name: "Last synced", Value: Text(peer.LastSyncedAt)},
		Field{Name: "Last error", Value: Text(peer.LastError)},
		Field{Name: "Created", Value: Text(peer.CreatedAt)},
		Field{Name: "Updated", Value: Text(peer.UpdatedAt)},
	)
}

// CheckToRecord renders the answer of a check: OK, status, version, message.
func (s SyncPeerSerializer) CheckToRecord(check domain.SyncCheckModel) Record {
	return Record{
		{Name: "OK", Value: YesNo(check.OK)},
		{Name: "Status", Value: Number(check.Status)},
		{Name: "Version", Value: Text(check.Version)},
		{Name: "Message", Value: check.Message},
	}
}

// InboundTokenToRecord renders a new inbound token by its prefix; the
// plain token is printed apart from it (see [SyncPeerController.InboundToken]).
func (s SyncPeerSerializer) InboundTokenToRecord(token domain.SyncInboundTokenModel) Record {
	return Record{{Name: "Token prefix", Value: token.TokenPrefix}}
}

// PushToRecord renders the answer of a push or a run: whether it was a dry
// run, how many pushes were queued, to which peers, a line per release and
// per skipped release.
func (s SyncPeerSerializer) PushToRecord(push domain.SyncPushModel) Record {
	peers := make([]string, len(push.Peers))
	for i, peer := range push.Peers {
		peers[i] = strconv.Itoa(peer)
	}
	record := Record{
		{Name: "Queued", Value: strconv.Itoa(push.Queued)},
		{Name: "Dry run", Value: YesNo(push.DryRun)},
		{Name: "Peers", Value: List(peers)},
		{Name: "Releases", Value: strconv.Itoa(len(push.Releases))},
	}
	for _, release := range push.Releases {
		record = append(record, Field{Name: "Release", Value: release.Software + " " + release.Version})
	}
	for _, skipped := range push.Skipped {
		record = append(record, Field{Name: "Skipped", Value: skipped})
	}
	return record
}
