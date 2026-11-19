package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// software adds the "software" commands.
func (r Router) software(root *cobra.Command) {
	serializer := SoftwareSerializer{}
	software := group(root, "catalog", "software", "The softwares of the catalog: games, demos, tools", "softwares", "sw")
	add(software, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the softwares", Args: ExactArgs(0),
		Example: "  warp-engine software list --platform c64\n  warp-engine software list --query ra --all -o json",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, SoftwareSorts)
		},
		Run: r.SoftwareController.Index,
	})
	add(software, "", Route{
		Use: "show SOFTWARE", Short: "Show a software (by ID or name)", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine software show blessing-of-ra --include releases,images",
		Flags: func(flags *pflag.FlagSet) {
			flags.StringSlice("include", nil, "embed related records: releases, images, externalLinks, extras")
		},
		Run: r.SoftwareController.Show,
	})
	add(software, "", Route{
		Use: "create", Short: "Create a software", Args: ExactArgs(0),
		Example: "  warp-engine software create --name pong --title Pong --platform love",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareController.Create,
	})
	add(software, "", Route{
		Use: "update SOFTWARE", Short: "Change the fields given as flags", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine software update pong --status released --highlighted\n  warp-engine software update pong --author \"\"   # clears the author",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareController.Update,
	})
	add(software, "", Route{
		Use: "delete SOFTWARE", Aliases: []string{"rm"}, Short: "Delete a software", Args: ExactArgs(1, "SOFTWARE"),
		Flags: DeclareYes, Run: r.SoftwareController.Destroy,
	})
	add(software, "", Route{
		Use: "scaffold", Short: "Create a software with its git repository and CI pipeline", Args: ExactArgs(0),
		Example: "  warp-engine software scaffold --name my-game --title \"My Game\" --platform tic80",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.ScaffoldFields()) },
		Run:     r.SoftwareController.Scaffold,
	})
}
