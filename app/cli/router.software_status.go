package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// softwareStatus adds the "status" commands.
func (r Router) softwareStatus(root *cobra.Command) {
	serializer := SoftwareStatusSerializer{}
	status := group(root, "catalog", "status", "The software statuses: the words a software's status takes", "statuses")
	add(status, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the software statuses", Args: ExactArgs(0),
		Example: "  warp-engine status list --sort key\n  warp-engine status list --default-visible=false",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, SoftwareStatusSorts)
		},
		Run: r.SoftwareStatusController.Index,
	})
	add(status, "", Route{
		Use: "show ID", Short: "Show a status", Args: ExactArgs(1, "ID"),
		Run: r.SoftwareStatusController.Show,
	})
	add(status, "", Route{
		Use: "create", Short: "Create a status", Args: ExactArgs(0),
		Example: "  warp-engine status create --key beta --label Beta --badge BETA --default-visible",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareStatusController.Create,
	})
	add(status, "", Route{
		Use: "update ID", Short: "Change the fields given as flags", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine status update 4 --position 2 --badge \"\"   # clears the badge",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareStatusController.Update,
	})
	add(status, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a status", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SoftwareStatusController.Destroy,
	})
}
