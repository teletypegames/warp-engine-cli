package cli

import (
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// BuildPlatformControllerInterface completes the flags that take the
// server's platforms, asset kinds and app variants in the shell. Its
// actions have cobra's completion signature.
type BuildPlatformControllerInterface interface {
	// Platforms completes a --platform.
	Platforms(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective)
	// PipelinePlatforms completes the --platform of a pipeline: the
	// platforms and [domain.PipelineUnknownPlatform].
	PipelinePlatforms(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective)
	// AssetKinds completes an asset's --kind.
	AssetKinds(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective)
	// Apps completes the --app of the --platform given on the same line.
	Apps(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective)
}

// BuildPlatformController is the [BuildPlatformControllerInterface] over a
// [domain.BuildPlatformServiceInterface]. When the server cannot be asked,
// it completes nothing.
type BuildPlatformController struct {
	// BuildPlatformService knows the server's platforms.
	BuildPlatformService domain.BuildPlatformServiceInterface
}

// Platforms returns the server's platforms that start with to_complete.
func (c BuildPlatformController) Platforms(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective) {
	return completions(c.BuildPlatformService.Platforms(), to_complete)
}

// PipelinePlatforms returns the server's platforms and "unknown" that
// start with to_complete.
func (c BuildPlatformController) PipelinePlatforms(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective) {
	platforms := c.BuildPlatformService.Platforms()
	if platforms != nil {
		platforms = append(platforms, domain.PipelineUnknownPlatform)
	}
	return completions(platforms, to_complete)
}

// AssetKinds returns the asset kinds of the server's platforms that start
// with to_complete.
func (c BuildPlatformController) AssetKinds(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective) {
	return completions(c.BuildPlatformService.AssetKinds(), to_complete)
}

// Apps returns the variants of the --platform on the command line that
// start with to_complete; nothing without a --platform.
func (c BuildPlatformController) Apps(cmd *cobra.Command, args []string, to_complete string) ([]string, cobra.ShellCompDirective) {
	platform, _ := cmd.Flags().GetString("platform")
	if platform == "" {
		return completions(nil, to_complete)
	}
	apps, _ := c.BuildPlatformService.Apps(platform)
	return completions(apps, to_complete)
}

// completions returns the values that start with to_complete, and tells
// the shell not to offer file names instead.
func completions(values []string, to_complete string) ([]string, cobra.ShellCompDirective) {
	matching := []string{}
	for _, value := range values {
		if strings.HasPrefix(value, to_complete) {
			matching = append(matching, value)
		}
	}
	return matching, cobra.ShellCompDirectiveNoFileComp
}
