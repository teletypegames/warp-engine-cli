package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Version is the program's version, set at build time:
//
//	go build -ldflags "-X git.teletypegames.org/engines/warp-engine-cli/app/cli.Version=1.2.0"
var Version = "dev"

// offline is the annotation of a command that runs without connecting to
// a server (login, profile, version, completion).
const offline = "offline"

// Route is one command of the tree: how it is called, its help and the
// controller action it runs. The router turns each into a cobra command.
type Route struct {
	// Use is the command's name and its arguments: "show SOFTWARE".
	Use string
	// Aliases are other names of the command ("ls" for "list").
	Aliases []string
	// Short is the one-line help.
	Short string
	// Long is the full help; empty uses Short.
	Long string
	// Example shows how it is called, indented two spaces.
	Example string
	// Args checks the positional arguments; nil takes none.
	Args cobra.PositionalArgs
	// Flags adds the command's flags.
	Flags func(flags *pflag.FlagSet)
	// Run is the controller action.
	Run func(cmd *cobra.Command, args []string) error
	// Offline marks a command that does not connect to a server first.
	Offline bool
}

// Router builds the command tree and binds every command to its controller
// action, the way an HTTP router binds paths. Each field is one controller.
type Router struct {
	// AuthController: login, logout, whoami, server, and the connection.
	AuthController AuthControllerInterface
	// ProfileController: profile list|show|use|remove.
	ProfileController ProfileControllerInterface
	// SoftwareController: software list|show|create|update|delete|scaffold.
	SoftwareController SoftwareControllerInterface
	// RawController: api METHOD PATH.
	RawController RawControllerInterface
	// ReleaseController: release list|show|create|update|delete.
	ReleaseController ReleaseControllerInterface
	// ReleaseAssetController: asset list|show|create|update|delete.
	ReleaseAssetController ReleaseAssetControllerInterface
	// ExternalLinkController: link list|show|create|update|delete.
	ExternalLinkController ExternalLinkControllerInterface
	// ForkLinkController: fork list|show|update|delete.
	ForkLinkController ForkLinkControllerInterface
	// SoftwareImageController: image list|show|create|upload|update|delete.
	SoftwareImageController SoftwareImageControllerInterface
	// ExtraController: extra show|update.
	ExtraController ExtraControllerInterface
	// SoftwareStatusController: status list|show|create|update|delete.
	SoftwareStatusController SoftwareStatusControllerInterface
	// SoftwareCategoryController: category list|show|create|update|delete.
	SoftwareCategoryController SoftwareCategoryControllerInterface
	// PlatformLinkController: platform-link list|show|create|update|delete.
	PlatformLinkController PlatformLinkControllerInterface
	// DownloadController: download list|show|summary|top-files|top-softwares|stats.
	DownloadController DownloadControllerInterface
	// PipelineController: pipeline list|show|update|sync|runs|trigger.
	PipelineController PipelineControllerInterface
	// FileController: file list|upload|rename|delete|mkdir.
	FileController FileControllerInterface
	// RuntimeController: runtime list|show|install.
	RuntimeController RuntimeControllerInterface
	// LocalComponentController: local components|enable|disable|remove|credentials|ci-setup|ci-teardown.
	LocalComponentController LocalComponentControllerInterface
	// LocalOperationController: local operation (added only with the
	// LocalComponentController).
	LocalOperationController LocalOperationControllerInterface
	// SyncPeerController: peer list|show|create|update|delete|check|inbound-token|push|run.
	SyncPeerController SyncPeerControllerInterface
	// SyncPairingController: pairing list|show|create|update|delete|token.
	SyncPairingController SyncPairingControllerInterface
	// SyncEventController: sync-event list|show|software|status.
	SyncEventController SyncEventControllerInterface
	// ApplicationTokenController: token list|show|create|update|delete|rotate.
	ApplicationTokenController ApplicationTokenControllerInterface
}

// Init returns the root command with every command under it.
func (r Router) Init() *cobra.Command {
	root := &cobra.Command{
		Use:           "warp-engine",
		Short:         "Manage a WarpEngine server from the command line",
		Long:          rootHelp,
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if isOffline(cmd) {
				return nil
			}
			return r.AuthController.Connect(cmd, args)
		},
	}
	root.SetVersionTemplate(fmt.Sprintf("warp-engine %s (%s, %s/%s)\n", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH))
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error { return &UsageError{Err: err} })
	flags := root.PersistentFlags()
	flags.StringP("profile", "p", "", "the profile to use (env WARP_ENGINE_PROFILE; default: the current one)")
	flags.StringP("server", "s", "", "a server URL instead of a profile (env WARP_ENGINE_URL)")
	flags.String("token", "", "an application token instead of the profile's (env WARP_ENGINE_TOKEN)")
	flags.StringP("output", "o", OutputTable, "output format: table, json or id")
	root.AddGroup(
		&cobra.Group{ID: "session", Title: "Signing in:"},
		&cobra.Group{ID: "catalog", Title: "The catalog:"},
		&cobra.Group{ID: "operations", Title: "Builds, files and statistics:"},
		&cobra.Group{ID: "sync", Title: "Sync with other engines:"},
		&cobra.Group{ID: "admin", Title: "Administration:"},
	)
	root.SetHelpCommandGroupID("admin")
	root.SetCompletionCommandGroupID("admin")

	r.session(root)
	if r.ProfileController != nil {
		r.profile(root)
	}
	if r.SoftwareController != nil {
		r.software(root)
	}
	if r.RawController != nil {
		r.raw(root)
	}
	if r.ReleaseController != nil {
		r.release(root)
	}
	if r.ReleaseAssetController != nil {
		r.releaseAsset(root)
	}
	if r.ExternalLinkController != nil {
		r.externalLink(root)
	}
	if r.ForkLinkController != nil {
		r.forkLink(root)
	}
	if r.SoftwareImageController != nil {
		r.softwareImage(root)
	}
	if r.ExtraController != nil {
		r.extra(root)
	}
	if r.SoftwareStatusController != nil {
		r.softwareStatus(root)
	}
	if r.SoftwareCategoryController != nil {
		r.softwareCategory(root)
	}
	if r.PlatformLinkController != nil {
		r.platformLink(root)
	}
	if r.DownloadController != nil {
		r.download(root)
	}
	if r.PipelineController != nil {
		r.pipeline(root)
	}
	if r.FileController != nil {
		r.file(root)
	}
	if r.RuntimeController != nil {
		r.runtime(root)
	}
	if r.LocalComponentController != nil {
		r.localComponent(root)
	}
	if r.SyncPeerController != nil {
		r.syncPeer(root)
	}
	if r.SyncPairingController != nil {
		r.syncPairing(root)
	}
	if r.SyncEventController != nil {
		r.syncEvent(root)
	}
	if r.ApplicationTokenController != nil {
		r.applicationToken(root)
	}
	return root
}

// group adds a command that only holds subcommands under parent, in the
// help group group_id ("" for none).
func group(parent *cobra.Command, group_id string, use string, short string, aliases ...string) *cobra.Command {
	command := &cobra.Command{Use: use, Short: short, Aliases: aliases, GroupID: group_id, Args: ExactArgs(0, "a subcommand"),
		RunE: func(cmd *cobra.Command, args []string) error { return cmd.Help() }}
	parent.AddCommand(command)
	return command
}

// add turns route into a command under parent, in the help group group_id.
func add(parent *cobra.Command, group_id string, route Route) *cobra.Command {
	args := route.Args
	if args == nil {
		args = ExactArgs(0)
	}
	command := &cobra.Command{
		Use: route.Use, Aliases: route.Aliases, Short: route.Short, Long: route.Long, Example: route.Example,
		Args: args, RunE: route.Run, GroupID: group_id,
	}
	if route.Long == "" {
		command.Long = route.Short
	}
	if route.Offline {
		command.Annotations = map[string]string{offline: "true"}
	}
	if route.Flags != nil {
		route.Flags(command.Flags())
	}
	parent.AddCommand(command)
	return command
}

// isOffline reports whether cmd runs without a connection: it is marked
// offline, or it is cobra's help or completion.
func isOffline(cmd *cobra.Command) bool {
	if cmd.Annotations[offline] == "true" {
		return true
	}
	for current := cmd; current != nil; current = current.Parent() {
		if current.Name() == "completion" || current.Name() == "help" || current.Name() == cobra.ShellCompRequestCmd {
			return true
		}
	}
	return false
}

// rootHelp is the long help of the root command.
const rootHelp = `warp-engine manages a WarpEngine server through its Admin API: the catalog
(softwares, releases, files, images, links), its vocabulary, the pipelines,
the artifact store, the runtimes, the downloads, sync with other engines, the
application tokens and the components of WarpEngine Local.

Start by signing in; the server becomes the current profile:

  warp-engine login https://games.example.org
  warp-engine software list

Every list takes --page/--per-page/--all, every command -o json for scripts.
The full manual: https://wiki.teletypegames.org/godoc/warp-engine-cli/latest/cmd/warp-engine/`

// loginHelp is the long help of "login".
const loginHelp = `Sign in to a WarpEngine server and save it as a profile, which becomes the
current one.

Without URL it signs in to WarpEngine Local on this machine
(http://127.0.0.1:7470) with its token file, and so does it for any server
that offers the token file. Any other server is signed in to with a device
sign-in: the command prints a code and a page, opens the page in the browser,
and waits until someone with admin rights approves it there.

With --token it saves a token you already have (an application token made in
the admin, or by "warp-engine token create"), after checking it.`
