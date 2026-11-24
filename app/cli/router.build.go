package cli

import (
	"math"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// build adds "release upload" and "release publish" under the "release"
// command, making it when there is none.
func (r Router) build(root *cobra.Command) {
	serializer := BuildSerializer{}
	release := subcommand(root, "release")
	if release == nil {
		release = group(root, "catalog", "release", "The releases (versions) of a software", "releases")
	}
	add(release, "", Route{
		Use: "upload FILE...", Short: "Upload files built outside the CI into a release's folder (application token, upload scope)",
		Long: "Upload files into the folder of a release, as the CI does: for the packages the CI\n" +
			"does not build (macOS, signed on a Mac). Every file name starts with NAME-VERSION;\n" +
			"the server checks the SHA-256 the command sends. Then \"release publish\" adds\n" +
			"them to the release.\n\n" +
			"The build endpoints take an application token with the upload scope: --token,\n" +
			"WARP_TOKEN or WARP_ENGINE_TOKEN; the server is --server, WARP_SERVER,\n" +
			"WARP_ENGINE_URL or the profile.",
		Example: "  WARP_SERVER=https://teletypegames.org WARP_TOKEN=we_xxx \\\n" +
			"    warp-engine release upload --name pong --version 1.2.0 pong-1.2.0-mac-x64.zip pong-1.2.0-mac-arm64.zip",
		Args:    RangeArgs(1, math.MaxInt, "FILE..."),
		Flags:   func(flags *pflag.FlagSet) { serializer.DeclareTarget(flags, false) },
		Run:     r.BuildController.Upload,
		Offline: true,
	})
	add(release, "", Route{
		Use: "publish", Short: "Publish a release with every file in its folder (application token, update scope)",
		Long: "Register the uploaded files of a release, as the CI's publish step does; publishing\n" +
			"a release again keeps its files and adds the new ones. The token needs the update\n" +
			"scope; the server and the token are found as for \"release upload\".",
		Example: "  warp-engine release publish --name pong --version 1.2.0 --platform ebitengine\n" +
			"  warp-engine release publish --name gio-demo --version 0.3.0 --platform golang --app gui",
		Flags:   func(flags *pflag.FlagSet) { serializer.DeclareTarget(flags, true) },
		Run:     r.BuildController.Publish,
		Offline: true,
	})
}

// subcommand returns the command named name under parent, or nil.
func subcommand(parent *cobra.Command, name string) *cobra.Command {
	for _, command := range parent.Commands() {
		if command.Name() == name {
			return command
		}
	}
	return nil
}
