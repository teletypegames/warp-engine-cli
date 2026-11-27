package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// softwareSource adds the "source" commands.
func (r Router) softwareSource(root *cobra.Command) {
	serializer := SoftwareSourceSerializer{}
	source := group(root, "catalog", "source", "The source repository of a software, one per software")
	add(source, "", Route{
		Use: "show SOFTWARE", Short: "Show the source of a software (by ID or name)", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine source show pong",
		Run:     r.SoftwareSourceController.Show,
	})
	add(source, "", Route{
		Use: "set SOFTWARE", Short: "Set the source of a software, replacing the one it had", Args: ExactArgs(1, "SOFTWARE"),
		Long: "Set the git repository of a software (by ID or name). The source is replaced as a\n" +
			"whole: --url is required, a missing --ref means the default branch.",
		Example: "  warp-engine source set pong --url https://git.example.org/games/pong --ref main",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareSourceController.Set,
	})
	add(source, "", Route{
		Use: "delete SOFTWARE", Aliases: []string{"rm"}, Short: "Remove the source of a software", Args: ExactArgs(1, "SOFTWARE"),
		Flags: DeclareYes, Run: r.SoftwareSourceController.Destroy,
	})
}
