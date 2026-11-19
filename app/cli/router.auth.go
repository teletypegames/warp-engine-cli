package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// session adds the sign-in commands.
func (r Router) session(root *cobra.Command) {
	add(root, "session", Route{
		Use: "login [URL]", Short: "Sign in to a server and save it as a profile", Offline: true,
		Long: loginHelp,
		Example: "  warp-engine login https://games.example.org\n" +
			"  warp-engine login                      # WarpEngine Local on this machine\n" +
			"  warp-engine login https://games.example.org --token we_xxx --name ci",
		Args: RangeArgs(0, 1, "[URL]"),
		Flags: func(flags *pflag.FlagSet) {
			flags.String("name", "", "the profile name (default: the host, or \"local\")")
			flags.Bool("local", false, "use the WarpEngine Local token file")
			flags.StringSlice("scopes", nil, "scopes to ask for (default admin:read,admin:write)")
			flags.Bool("no-browser", false, "do not open the approval page in the browser")
		},
		Run: r.AuthController.Login,
	})
	add(root, "session", Route{
		Use: "logout [PROFILE]", Short: "Revoke the token and forget the profile", Offline: true,
		Args: RangeArgs(0, 1, "[PROFILE]"),
		Flags: func(flags *pflag.FlagSet) {
			flags.Bool("keep-token", false, "forget the profile but do not revoke the token")
		},
		Run: r.AuthController.Logout,
	})
	add(root, "session", Route{Use: "whoami", Short: "Show who the server takes the token for", Args: ExactArgs(0), Run: r.AuthController.Whoami})
	add(root, "session", Route{Use: "server", Short: "Show what the server says about itself", Args: ExactArgs(0), Run: r.AuthController.Server})

	add(root, "admin", Route{
		Use: "version", Short: "Show the version of this program", Offline: true, Args: ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "warp-engine %s (%s, %s/%s)\n", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
			return nil
		},
	})
}
