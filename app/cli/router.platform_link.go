package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// platformLink adds the "platform-link" commands.
func (r Router) platformLink(root *cobra.Command) {
	serializer := PlatformLinkSerializer{}
	link := group(root, "catalog", "platform-link", "The links shown with every software of a platform", "platform-links")
	add(link, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the platform links", Args: ExactArgs(0),
		Example: "  warp-engine platform-link list --platform c64",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, PlatformLinkSorts)
		},
		Run: r.PlatformLinkController.Index,
	})
	add(link, "", Route{
		Use: "show ID", Short: "Show a platform link", Args: ExactArgs(1, "ID"),
		Run: r.PlatformLinkController.Show,
	})
	add(link, "", Route{
		Use: "create", Short: "Create a platform link", Args: ExactArgs(0),
		Example: "  warp-engine platform-link create --platform c64 --name VICE --url https://vice-emu.sourceforge.io/",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.PlatformLinkController.Create,
	})
	add(link, "", Route{
		Use: "update ID", Short: "Change the fields given as flags", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine platform-link update 3 --position 0",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.PlatformLinkController.Update,
	})
	add(link, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a platform link", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.PlatformLinkController.Destroy,
	})
	if r.BuildPlatformController != nil {
		r.completeFlag(link, "platform", r.BuildPlatformController.Platforms)
	}
}
