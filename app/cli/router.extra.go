package cli

import "github.com/spf13/cobra"

// extra adds the "extra" commands.
func (r Router) extra(root *cobra.Command) {
	serializer := ExtraSerializer{}
	extra := group(root, "catalog", "extra", "The host's additional resources of a software (hosted engines only)", "extras")
	add(extra, "", Route{
		Use: "show SOFTWARE ASSOCIATION", Short: "Show what a host resource holds for a software", Args: ExactArgs(2, "SOFTWARE", "ASSOCIATION"),
		Long: "Show what one of the host's additional resources holds for a software. Which\n" +
			"resources there are depends on the host (\"software show SOFTWARE --include extras\"\n" +
			"shows them all); an engine without any answers 404.",
		Example: "  warp-engine extra show pong seo_metadata -o json",
		Run:     r.ExtraController.Show,
	})
	add(extra, "", Route{
		Use: "update SOFTWARE ASSOCIATION", Short: "Replace what a host resource holds for a software", Args: ExactArgs(2, "SOFTWARE", "ASSOCIATION"),
		Long: "Replace what one of the host's additional resources holds for a software with\n" +
			"--value: a JSON object, or an array of objects when the resource holds many\n" +
			"(an entry with \"id\" and \"_destroy\": true removes one). --data takes the whole\n" +
			"body, {\"value\": …}.",
		Example: "  warp-engine extra update pong seo_metadata --value '{\"title\":\"Pong\"}'\n  warp-engine extra update pong seo_metadata --value @seo.json",
		Flags:   serializer.DeclareFlags,
		Run:     r.ExtraController.Update,
	})
}
