package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// applicationToken adds the "token" commands.
func (r Router) applicationToken(root *cobra.Command) {
	serializer := ApplicationTokenSerializer{}
	token := group(root, "admin", "token", "Application tokens: what programs sign in with", "tokens")
	add(token, "", Route{
		Use: "list", Aliases: []string{"ls"}, Short: "List the application tokens", Args: ExactArgs(0),
		Example: "  warp-engine token list --kind publishing\n  warp-engine token list --state expired",
		Flags: func(flags *pflag.FlagSet) {
			DeclareFilters(flags, serializer.Filters())
			DeclareListFlags(flags, ApplicationTokenSorts)
		},
		Run: r.ApplicationTokenController.Index,
	})
	add(token, "", Route{
		Use: "show ID", Short: "Show an application token", Args: ExactArgs(1, "ID"),
		Run: r.ApplicationTokenController.Show,
	})
	add(token, "", Route{
		Use: "create", Short: "Issue a token (the token is shown once)", Args: ExactArgs(0),
		Example: "  warp-engine token create --name ci --scopes admin:read,admin:write --expires-at 2027-06-30T00:00:00Z\n" +
			"  warp-engine token create --name reader --scopes admin:read -o json | jq -r .plainToken",
		Flags: func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.CreateFields()) },
		Run:   r.ApplicationTokenController.Create,
	})
	add(token, "", Route{
		Use: "update ID", Short: "Rename a token, or change its scopes or expiry", Args: ExactArgs(1, "ID"),
		Example: "  warp-engine token update 12 --scopes admin:read\n  warp-engine token update 12 --expires-at \"\"   # never expires",
		Flags:   func(flags *pflag.FlagSet) { DeclareFields(flags, serializer.Fields()) },
		Run:     r.ApplicationTokenController.Update,
	})
	add(token, "", Route{
		Use: "delete ID", Aliases: []string{"rm", "revoke"}, Short: "Revoke a token", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.ApplicationTokenController.Destroy,
	})
	add(token, "", Route{
		Use: "rotate ID", Short: "Replace a token with a new one (shown once) and push it to its CI repositories", Args: ExactArgs(1, "ID"),
		Flags: DeclareYes, Run: r.ApplicationTokenController.Rotate,
	})
}
