package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// releaseAsset adds the "asset" commands.
func (r Router) releaseAsset(root *cobra.Command) {
	serializer := ReleaseAssetSerializer{}
	asset := group(root, "catalog", "asset", "The files of a release in the artifact store", "assets")
	add(asset, "", Route{
		Use: "list RELEASE", Aliases: []string{"ls"}, Short: "List the files of a release", Args: ExactArgs(1, "RELEASE"),
		Example: "  warp-engine asset list 42",
		Run:     r.ReleaseAssetController.Index,
	})
	add(asset, "", Route{
		Use: "show ID", Short: "Show a release file", Args: ExactArgs(1, "ID"),
		Run: r.ReleaseAssetController.Show,
	})
	add(asset, "", Route{
		Use: "create RELEASE", Short: "Add a file of the artifact store to a release", Args: ExactArgs(1, "RELEASE"),
		Example: "  warp-engine asset create 42 --kind html --path /pong/1.2.0/html/index.html",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.ReleaseAssetController.Create,
	})
	add(asset, "", Route{
		Use: "update ID", Short: "Change the kind or path of a release file", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine asset update 7 --path /pong/1.2.0/pong.love",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.ReleaseAssetController.Update,
	})
	add(asset, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Delete a release file", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.ReleaseAssetController.Destroy,
	})
}
