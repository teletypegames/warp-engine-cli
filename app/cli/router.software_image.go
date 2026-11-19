package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// softwareImage adds the "image" commands.
func (r Router) softwareImage(root *cobra.Command) {
	serializer := SoftwareImageSerializer{}
	image := group(root, "catalog", "image", "The images (screenshots, covers) of a software", "images")
	add(image, "", Route{
		Use: "list SOFTWARE", Aliases: []string{"ls"}, Short: "List the images of a software in order (by ID or name)", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine image list pong",
		Run:     r.SoftwareImageController.Index,
	})
	add(image, "", Route{
		Use: "show ID", Short: "Show a software image", Args: ExactArgs(1, "ID"),
		Run: r.SoftwareImageController.Show,
	})
	add(image, "", Route{
		Use: "create SOFTWARE", Short: "Bind an image of the host's image library to a software", Args: ExactArgs(1, "SOFTWARE"),
		Example: "  warp-engine image create pong --image-id 118 --position 2",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.CreateFields()) },
		Run:     r.SoftwareImageController.Create,
	})
	add(image, "", Route{
		Use: "upload SOFTWARE FILE", Short: "Upload an image file into the library and bind it to a software", Args: ExactArgs(2, "SOFTWARE", "FILE"),
		Example: "  warp-engine image upload pong cover.png --default\n  warp-engine image upload pong shot-2.png --position 2",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareImageController.Upload,
	})
	add(image, "", Route{
		Use: "update ID", Short: "Make an image the default, or move it", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine image update 31 --default\n  warp-engine image update 31 --position 0",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.SoftwareImageController.Update,
	})
	add(image, "", Route{
		Use: "delete ID", Aliases: []string{"rm"}, Short: "Unbind an image from its software (the library keeps it)", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.SoftwareImageController.Destroy,
	})
}
