package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// softwareLink adds the "link" commands.
func (r Router) softwareLink(root *cobra.Command) {
	serializer := SoftwareLinkSerializer{}
	link := group(root, "catalog", "link", "The links of a software (its source is \"source\")", "links")
	add(link, "", Route{
		Use: "list SOFTWARE", Aliases: []string{"ls"}, Short: "List the links of a software (by ID or name)", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine link list pong",
		Run:     r.SoftwareLinkController.Index,
	})
	add(link, "", Route{
		Use: "show ID", Short: "Show a link", Args: ExactArgs(1, "ID"),
		Run: r.SoftwareLinkController.Show,
	})
	add(link, "", Route{
		Use: "create SOFTWARE", Short: "Add a link to a software", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine link create pong --label Store --url https://store.example.org/pong",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareLinkController.Create,
	})
	add(link, "", Route{
		Use: "update ID", Short: "Change the label or URL of a link", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine link update 5 --label \"Store page\"",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareLinkController.Update,
	})
	add(link, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a link", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SoftwareLinkController.Destroy,
	})
}
