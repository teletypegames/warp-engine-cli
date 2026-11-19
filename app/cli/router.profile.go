package cli

import "github.com/spf13/cobra"

// profile adds the "profile" commands.
func (r Router) profile(root *cobra.Command) {
	profile := group(root, "session", "profile", "The saved servers and tokens", "profiles")
	add(profile, "", Route{Use: "list", Aliases: []string{"ls"}, Short: "List the profiles", Offline: true, Args: ExactArgs(0), Run: r.ProfileController.Index})
	add(profile, "", Route{Use: "show [NAME]", Short: "Show a profile (the current one without NAME)", Offline: true, Args: RangeArgs(0, 1, "[NAME]"), Run: r.ProfileController.Show})
	add(profile, "", Route{Use: "use NAME", Short: "Make a profile the current one", Offline: true, Args: ExactArgs(1, "NAME"), Run: r.ProfileController.Use})
	add(profile, "", Route{Use: "remove NAME", Aliases: []string{"rm"}, Short: "Forget a profile without revoking its token", Offline: true, Args: ExactArgs(1, "NAME"), Run: r.ProfileController.Destroy})
}
