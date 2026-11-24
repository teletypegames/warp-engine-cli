package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// makefile adds "check makefile".
func (r Router) makefile(root *cobra.Command) {
	check := subcommand(root, "check")
	if check == nil {
		check = group(root, "operations", "check", "Check a project against the WarpEngine conventions")
	}
	add(check, "", Route{
		Use: "makefile [DIR]", Short: "Check that a project's Makefile has the targets of the shared contract",
		Long: "Check the Makefile of the project in DIR (the current directory without it)\n" +
			"against the contract of the WarpEngine templates: build, run, check and clean\n" +
			"always; web when the platform's releases have a web build (html); mac-release\n" +
			"when the platform has kinds the CI does not build. Every target is tried with\n" +
			"make -n, which builds nothing. The platform comes from the CI marker\n" +
			"(.woodpecker.yaml \"platform:\" and \"app:\") or --platform/--app, its kinds from\n" +
			"the server (--server, WARP_SERVER, WARP_ENGINE_URL or the profile); without a\n" +
			"server only the four targets every project has are checked.\n\n" +
			"Exits with 1 when a target is missing or make -n fails on it.",
		Example: "  warp-engine check makefile\n  warp-engine check makefile ../pong --server https://teletypegames.org\n" +
			"  warp-engine check makefile . --platform golang --app gui",
		Args: RangeArgs(0, 1, "[DIR]"),
		Flags: func(flags *pflag.FlagSet) {
			flags.String("platform", "", "the platform, instead of the CI marker's")
			flags.String("app", "", "the platform variant, instead of the CI marker's")
		},
		Run:     r.MakefileController.Check,
		Offline: true,
	})
}
