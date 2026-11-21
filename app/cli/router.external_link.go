package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// externalLink adds the "link" commands.
func (r Router) externalLink(root *cobra.Command) {
	serializer := ExternalLinkSerializer{}
	link := group(root, "catalog", "link", "The external links of a software", "links")
	add(link, "", Route{
		Use: "list SOFTWARE", Aliases: []string{"ls"}, Short: "List the links of a software (by ID or name)", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine link list pong",
		Run:     r.ExternalLinkController.Index,
	})
	add(link, "", Route{
		Use: "show ID", Short: "Show a link", Args: ExactArgs(1, "ID"),
		Run: r.ExternalLinkController.Show,
	})
	add(link, "", Route{
		Use: "create SOFTWARE", Short: "Add a link to a software", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine link create pong --label Source --url https://git.example.org/games/pong --kind source --ref main",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.ExternalLinkController.Create,
	})
	add(link, "", Route{
		Use: "update ID", Short: "Change the label, URL, kind or ref of a link", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine link update 5 --label \"Source code\"",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.ExternalLinkController.Update,
	})
	add(link, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a link", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.ExternalLinkController.Destroy,
	})
}
