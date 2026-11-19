package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// softwareCategory adds the "category" commands.
func (r Router) softwareCategory(root *cobra.Command) {
	serializer := SoftwareCategorySerializer{}
	category := group(root, "catalog", "category", "The software categories: the words a software's category takes", "categories")
	add(category, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the software categories", Args: ExactArgs(0),
		Example: "  warp-engine category list --needs-review",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, SoftwareCategorySorts)
		},
		Run: r.SoftwareCategoryController.Index,
	})
	add(category, "", Route{
		Use: "show ID", Short: "Show a category", Args: ExactArgs(1, "ID"),
		Run: r.SoftwareCategoryController.Show,
	})
	add(category, "", Route{
		Use: "create", Short: "Create a category", Args: ExactArgs(0),
		Example: "  warp-engine category create --key tool --label Tool",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareCategoryController.Create,
	})
	add(category, "", Route{
		Use: "update ID", Short: "Change the fields given as flags", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine category update 9 --label Utility --needs-review=false   # approves it",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareCategoryController.Update,
	})
	add(category, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a category", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SoftwareCategoryController.Destroy,
	})
}
