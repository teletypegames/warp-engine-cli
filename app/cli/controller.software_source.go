package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SoftwareSourceControllerInterface is the "source" commands: the source
// repository of a software.
type SoftwareSourceControllerInterface interface {
	// Show shows the source of a software.
	Show(cmd *cobra.Command, args []string) error
	// Set sets the source of a software.
	Set(cmd *cobra.Command, args []string) error
	// Destroy removes the source of a software.
	Destroy(cmd *cobra.Command, args []string) error
}

// SoftwareSourceController is the [SoftwareSourceControllerInterface] over a
// [domain.SoftwareSourceServiceInterface].
type SoftwareSourceController struct {
	// SoftwareSourceService manages the sources.
	SoftwareSourceService domain.SoftwareSourceServiceInterface
}

// Show prints the source of SOFTWARE (an ID or a name).
func (c SoftwareSourceController) Show(cmd *cobra.Command, args []string) error {
	source, err := c.SoftwareSourceService.Show(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, source, SoftwareSourceSerializer{}.ModelToRecord(*source))
}

// Set sets the source of SOFTWARE from --url and --ref and prints it.
func (c SoftwareSourceController) Set(cmd *cobra.Command, args []string) error {
	serializer := SoftwareSourceSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	source, err := c.SoftwareSourceService.Set(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, source, serializer.ModelToRecord(*source))
}

// Destroy removes the source of SOFTWARE after asking (not with --yes).
func (c SoftwareSourceController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Remove the source of "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.SoftwareSourceService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Source of %s removed.", args[0])
	return nil
}
