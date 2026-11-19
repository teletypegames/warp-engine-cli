package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// release adds the "release" commands.
func (r Router) release(root *cobra.Command) {
	serializer := ReleaseSerializer{}
	release := group(root, "catalog", "release", "The releases (versions) of a software", "releases")
	add(release, "", Route{
		Use: "list SOFTWARE", Aliases: []string{"ls"}, Short: "List the releases of a software (by ID or name)", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine release list pong\n  warp-engine release list pong --channel beta --sort version\n  warp-engine release list pong --commit-sha 3fb3ed8 -o id",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, ReleaseSorts)
		},
		Run: r.ReleaseController.Index,
	})
	add(release, "", Route{
		Use: "show ID", Short: "Show a release with its files", Args: ExactArgs(1, "ID"),
		Run: r.ReleaseController.Show,
	})
	add(release, "", Route{
		Use: "create SOFTWARE", Short: "Add a release to a software; the version sets the channel", Args: ExactArgs(1, "SOFTWARE"),
		Long: "Add a release to a software. The channel follows from the version: dev-… is dev,\n" +
			"a pre-release such as 1.3.0-beta.1 is beta, anything else stable. A deleted\n" +
			"release of the same version comes back.",
		Example: "  warp-engine release create pong --version 1.2.0\n  warp-engine release create pong --version dev-3fb3ed8",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.ReleaseController.Create,
	})
	add(release, "", Route{
		Use: "update ID", Short: "Change the version of a release", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine release update 42 --version 1.2.1",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.ReleaseController.Update,
	})
	add(release, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a release with its files", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.ReleaseController.Destroy,
	})
}
