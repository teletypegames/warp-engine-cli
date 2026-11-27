package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// raw adds the "api" command.
func (r Router) raw(root *cobra.Command) {
	add(root, "admin", Route{
		Use: "api METHOD PATH", Short: "Send any request to the Admin API and print the answer",
		Long: `Send any request to the server and print its JSON answer as it came.

PATH is relative to the Admin API (/softwares means /api/admin/v2/softwares);
a path under /api/ or /.well-known/ is taken from the server's root. The
request is signed in like every other command. For what has no command of its
own, and to see the raw answers; the API's contract is GET /openapi.json.`,
		Example: "  warp-engine api GET /softwares --query platform=c64 --query perPage=5\n" +
			"  warp-engine api PATCH /softwares/12 --data '{\"highlighted\":true}'\n" +
			"  warp-engine api GET /openapi.json > openapi.json",
		Args: ExactArgs(2, "METHOD", "PATH"),
		Flags: func(flags *pflag.FlagSet) {
			flags.StringArray("query", nil, "a query parameter key=value; repeatable")
			flags.String(DataFlag, "", "the JSON body, or @file.json (@- reads standard input)")
		},
		Run: r.RawController.Send,
	})
}
