package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// ReleaseAssetControllerInterface is the "asset" commands: the files of a
// release.
type ReleaseAssetControllerInterface interface {
	// Index lists the assets of a release.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one asset.
	Show(cmd *cobra.Command, args []string) error
	// Create adds an asset to a release.
	Create(cmd *cobra.Command, args []string) error
	// Update changes an asset.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes an asset.
	Destroy(cmd *cobra.Command, args []string) error
}

// ReleaseAssetController is the [ReleaseAssetControllerInterface] over a
// [domain.ReleaseAssetServiceInterface].
type ReleaseAssetController struct {
	// ReleaseAssetService manages the release assets.
	ReleaseAssetService domain.ReleaseAssetServiceInterface
}

// Index prints the assets of the release RELEASE.
func (c ReleaseAssetController) Index(cmd *cobra.Command, args []string) error {
	assets, err := c.ReleaseAssetService.All(args[0])
	if err != nil {
		return err
	}
	return PrintList(cmd, assets, ReleaseAssetSerializer{}.ModelToTable(assets))
}

// Show prints the asset ID.
func (c ReleaseAssetController) Show(cmd *cobra.Command, args []string) error {
	asset, err := c.ReleaseAssetService.Show(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, asset, ReleaseAssetSerializer{}.ModelToRecord(*asset))
}

// Create adds an asset from --kind and --path to RELEASE and prints it.
func (c ReleaseAssetController) Create(cmd *cobra.Command, args []string) error {
	serializer := ReleaseAssetSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	asset, err := c.ReleaseAssetService.Create(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, asset, serializer.ModelToRecord(*asset))
}

// Update changes the fields given as flags on the asset ID and prints it.
func (c ReleaseAssetController) Update(cmd *cobra.Command, args []string) error {
	serializer := ReleaseAssetSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	asset, err := c.ReleaseAssetService.Update(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, asset, serializer.ModelToRecord(*asset))
}

// Destroy deletes the asset ID after asking (not with --yes).
func (c ReleaseAssetController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Delete the asset "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.ReleaseAssetService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Asset %s deleted.", args[0])
	return nil
}
