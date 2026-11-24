package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// BuildControllerInterface is "release upload" and "release publish": a
// release made outside the CI, through the build endpoints.
type BuildControllerInterface interface {
	// Upload sends files into a release's folder.
	Upload(cmd *cobra.Command, args []string) error
	// Publish registers a release with the files in its folder.
	Publish(cmd *cobra.Command, args []string) error
}

// BuildController is the [BuildControllerInterface] over a
// [domain.BuildServiceInterface]. Its commands connect themselves: they
// take the server and the token from WARP_SERVER and WARP_TOKEN too (see
// [ConnectBuild]).
type BuildController struct {
	// BuildService releases through the build endpoints.
	BuildService domain.BuildServiceInterface
	// AuthService connects.
	AuthService domain.AuthServiceInterface
}

// Upload uploads FILE… into the release --name --version and prints what
// the server stored.
func (c BuildController) Upload(cmd *cobra.Command, args []string) error {
	if err := ConnectBuild(cmd, c.AuthService); err != nil {
		return err
	}
	serializer := BuildSerializer{}
	uploads, err := c.BuildService.Upload(serializer.FlagsToTarget(cmd.Flags()), args)
	if len(uploads) > 0 {
		if print_err := PrintList(cmd, uploads, serializer.UploadsToTable(uploads)); print_err != nil && err == nil {
			err = print_err
		}
	}
	return err
}

// Publish publishes the release --name --version of --platform (--app,
// --commit) and prints the answer.
func (c BuildController) Publish(cmd *cobra.Command, args []string) error {
	if err := ConnectBuild(cmd, c.AuthService); err != nil {
		return err
	}
	serializer := BuildSerializer{}
	published, err := c.BuildService.Publish(serializer.FlagsToTarget(cmd.Flags()))
	if err != nil {
		return err
	}
	return PrintRecord(cmd, published, serializer.PublishToRecord(*published))
}

// ConnectBuild connects a command of the build endpoints: the server is
// --server, else WARP_SERVER, else WARP_ENGINE_URL; the token --token,
// else WARP_TOKEN, else WARP_ENGINE_TOKEN; without a server, the profile
// (--profile, WARP_ENGINE_PROFILE or the current one), as every command.
// WARP_SERVER and WARP_TOKEN are the names the Makefiles of the WarpEngine
// templates use.
func ConnectBuild(cmd *cobra.Command, auth domain.AuthServiceInterface) error {
	selection := domain.Selection{
		Profile: flagOrEnv(cmd, "profile", "WARP_ENGINE_PROFILE"),
		Server:  firstOf(flagOrEnv(cmd, "server", "WARP_SERVER"), flagOrEnv(cmd, "server", "WARP_ENGINE_URL")),
		Token:   firstOf(flagOrEnv(cmd, "token", "WARP_TOKEN"), flagOrEnv(cmd, "token", "WARP_ENGINE_TOKEN")),
	}
	_, err := auth.Connect(selection)
	return err
}

// firstOf returns the first non-empty value.
func firstOf(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
