package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// ProfileControllerInterface is the profile commands: the saved servers.
type ProfileControllerInterface interface {
	// Index lists the profiles.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one profile.
	Show(cmd *cobra.Command, args []string) error
	// Use makes a profile the current one.
	Use(cmd *cobra.Command, args []string) error
	// Destroy forgets a profile without signing out.
	Destroy(cmd *cobra.Command, args []string) error
}

// ProfileController is the [ProfileControllerInterface] over a
// [domain.ProfileServiceInterface].
type ProfileController struct {
	// ProfileService keeps the profiles.
	ProfileService domain.ProfileServiceInterface
}

// Index prints every profile, the current one marked. The JSON output
// leaves the tokens out.
func (c ProfileController) Index(cmd *cobra.Command, args []string) error {
	profiles, current, err := c.ProfileService.All()
	if err != nil {
		return err
	}
	serializer := ProfileSerializer{}
	return PrintList(cmd, serializer.ModelToResponseMany(*profiles), serializer.ModelToTable(*profiles, current))
}

// Show prints the profile NAME, or the current one without it.
func (c ProfileController) Show(cmd *cobra.Command, args []string) error {
	name := ""
	if len(args) == 1 {
		name = args[0]
	}
	profile, err := c.ProfileService.Show(name)
	if err != nil {
		return err
	}
	_, current, err := c.ProfileService.All()
	if err != nil {
		return err
	}
	serializer := ProfileSerializer{}
	return PrintRecord(cmd, serializer.ModelToResponse(*profile), serializer.ModelToRecord(*profile, current))
}

// Use makes NAME the current profile.
func (c ProfileController) Use(cmd *cobra.Command, args []string) error {
	if err := c.ProfileService.Use(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Profile %q is now current.", args[0])
	return nil
}

// Destroy forgets the profile NAME; its token stays valid on the server
// (use logout to revoke it).
func (c ProfileController) Destroy(cmd *cobra.Command, args []string) error {
	if err := c.ProfileService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Profile %q removed (its token was not revoked).", args[0])
	return nil
}
