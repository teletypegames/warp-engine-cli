package cli

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// DefaultOperationTimeout is how long the commands that start a WarpEngine
// Local operation wait for it by default (--timeout).
const DefaultOperationTimeout = 10 * time.Minute

// localComponent adds the "local" commands, "local operation" among them
// when the router has a LocalOperationController.
func (r Router) localComponent(root *cobra.Command) {
	local := group(root, "admin", "local", "WarpEngine Local: its git server, its CI, a software's local CI")
	local.Long = "The components WarpEngine Local runs next to the engine on this machine (a git\n" +
		"server, a CI) and the local CI of a software. Only WarpEngine Local has them:\n" +
		"any other server answers 404 (exit code 4).\n\n" +
		"Enabling, disabling and removing run in the background on the engine; the\n" +
		"commands wait for them and print their steps, --no-wait returns at once."
	waiting := func(flags *pflag.FlagSet) {
		flags.Bool("no-wait", false, "return the queued operation at once (follow it with \"local operation ID --wait\")")
		flags.Duration("timeout", DefaultOperationTimeout, "how long to wait for the operation")
	}
	add(local, "", Route{
		Use: "components", Short: "List the components with their state", Args: ExactArgs(0),
		Run: r.LocalComponentController.Index,
	})
	add(local, "", Route{
		Use: "enable KEY", Short: "Start a component and what it needs", Args: ExactArgs(1, "KEY"),
		Example: "  warp-engine local enable woodpecker\n  warp-engine local enable gitea --no-wait",
		Flags:   waiting, Run: r.LocalComponentController.Enable,
	})
	add(local, "", Route{
		Use: "disable KEY", Short: "Stop a component and what needs it (its data stays)", Args: ExactArgs(1, "KEY"),
		Flags: waiting, Run: r.LocalComponentController.Disable,
	})
	add(local, "", Route{
		Use: "remove KEY", Aliases: []string{"rm"}, Short: "Remove a component's containers (with --purge its data too)", Args: ExactArgs(1, "KEY"),
		Example: "  warp-engine local remove woodpecker --purge --yes",
		Flags: func(flags *pflag.FlagSet) {
			waiting(flags)
			DeclareYes(flags)
			flags.Bool("purge", false, "remove its data and generated sign-ins too")
		},
		Run: r.LocalComponentController.Destroy,
	})
	add(local, "", Route{
		Use: "credentials KEY", Short: "Show a component's generated sign-in (prints the password)", Args: ExactArgs(1, "KEY"),
		Run: r.LocalComponentController.Credentials,
	})
	add(local, "", Route{
		Use: "ci-setup SOFTWARE", Short: "Set a software up on the local git server and CI", Args: ExactArgs(1, "SOFTWARE"),
		Long: "Set a software (ID or name) up on the local git server and CI, starting them when\n" +
			"they are off. The result holds the push URL and the git remote to add.",
		Example: "  warp-engine local ci-setup my-game",
		Flags:   waiting, Run: r.LocalComponentController.SetupCI,
	})
	add(local, "", Route{
		Use: "ci-teardown SOFTWARE", Short: "Stop the local CI of a software (its repository stays)", Args: ExactArgs(1, "SOFTWARE"),
		Flags: DeclareYes, Run: r.LocalComponentController.TeardownCI,
	})
	if r.LocalOperationController != nil {
		r.localOperation(local)
	}
}
