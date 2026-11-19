package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SoftwareImageControllerInterface is the "image" commands: the pictures
// of a software.
type SoftwareImageControllerInterface interface {
	// Index lists the images of a software.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one software image.
	Show(cmd *cobra.Command, args []string) error
	// Create binds an image of the library to a software.
	Create(cmd *cobra.Command, args []string) error
	// Upload uploads a file into the library and binds it to a software.
	Upload(cmd *cobra.Command, args []string) error
	// Update makes an image the default or moves it.
	Update(cmd *cobra.Command, args []string) error
	// Destroy unbinds an image from its software.
	Destroy(cmd *cobra.Command, args []string) error
}

// SoftwareImageController is the [SoftwareImageControllerInterface] over a
// [domain.SoftwareImageServiceInterface].
type SoftwareImageController struct {
	// SoftwareImageService manages the software images.
	SoftwareImageService domain.SoftwareImageServiceInterface
}

// Index prints the images of SOFTWARE (an ID or a name) in order.
func (c SoftwareImageController) Index(cmd *cobra.Command, args []string) error {
	images, err := c.SoftwareImageService.All(args[0])
	if err != nil {
		return err
	}
	return PrintList(cmd, images, SoftwareImageSerializer{}.ModelToTable(images))
}

// Show prints the software image ID.
func (c SoftwareImageController) Show(cmd *cobra.Command, args []string) error {
	image, err := c.SoftwareImageService.Show(args[0])
	if err != nil {
		return err
	}
	return PrintRecord(cmd, image, SoftwareImageSerializer{}.ModelToRecord(*image))
}

// Create binds the library image --image-id to SOFTWARE (with --default
// and --position) and prints it.
func (c SoftwareImageController) Create(cmd *cobra.Command, args []string) error {
	serializer := SoftwareImageSerializer{}
	attributes, err := serializer.CreateToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	image, err := c.SoftwareImageService.Create(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, image, serializer.ModelToRecord(*image))
}

// Upload sends FILE into the image library, binds it to SOFTWARE (with
// --default and --position) and prints it.
func (c SoftwareImageController) Upload(cmd *cobra.Command, args []string) error {
	serializer := SoftwareImageSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	image, err := c.SoftwareImageService.Upload(args[0], args[1], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, image, serializer.ModelToRecord(*image))
}

// Update changes --default or --position of the software image ID and prints it.
func (c SoftwareImageController) Update(cmd *cobra.Command, args []string) error {
	serializer := SoftwareImageSerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	image, err := c.SoftwareImageService.Update(args[0], attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, image, serializer.ModelToRecord(*image))
}

// Destroy unbinds the software image ID after asking (not with --yes); the
// picture stays in the host's image library.
func (c SoftwareImageController) Destroy(cmd *cobra.Command, args []string) error {
	ok, err := Confirm(cmd, "Unbind the image "+args[0]+" from its software?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.SoftwareImageService.Destroy(args[0]); err != nil {
		return err
	}
	PrintDone(cmd, "Image %s unbound.", args[0])
	return nil
}
