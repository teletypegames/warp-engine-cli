package cli

import (
	"strconv"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/pflag"
)

// SoftwareImageSerializer turns the flags of the image commands into the
// API's terms and the software images into tables and records.
type SoftwareImageSerializer struct{}

// Fields are the field flags of "image update" and "image upload": whether
// the image is the default and where it stands.
func (s SoftwareImageSerializer) Fields() []FieldFlag {
	return []FieldFlag{
		{Flag: "default", Key: "isDefault", Kind: KindBool, Usage: "make it the image the stores show first"},
		{Flag: "position", Key: "position", Kind: KindInt, Usage: "its place in the order, from 0"},
	}
}

// CreateFields are the field flags of "image create": [SoftwareImageSerializer.Fields]
// and the image of the library to bind.
func (s SoftwareImageSerializer) CreateFields() []FieldFlag {
	return append([]FieldFlag{
		{Flag: "image-id", Key: "imageId", Kind: KindInt, Usage: "an image of the host's image library"},
	}, s.Fields()...)
}

// RequestToAttributes reads the field flags of update and upload.
func (s SoftwareImageSerializer) RequestToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.Fields())
}

// CreateToAttributes reads the field flags of create.
func (s SoftwareImageSerializer) CreateToAttributes(flags *pflag.FlagSet) (domain.Attributes, error) {
	return FieldsToAttributes(flags, s.CreateFields())
}

// ModelToTable renders software images as a table: ID, library image,
// position, default, origin, URL.
func (s SoftwareImageSerializer) ModelToTable(images []domain.SoftwareImageModel) Table {
	table := Table{Headers: []string{"ID", "Image", "Position", "Default", "Origin", "URL"}}
	for _, image := range images {
		table.Rows = append(table.Rows, []string{
			strconv.Itoa(image.ID), strconv.Itoa(image.ImageID), strconv.Itoa(image.Position),
			YesNo(image.IsDefault), image.Origin, Text(image.URL),
		})
	}
	return table
}

// ModelToRecord renders one software image with every field.
func (s SoftwareImageSerializer) ModelToRecord(image domain.SoftwareImageModel) Record {
	return Record{
		{Name: "ID", Value: strconv.Itoa(image.ID)},
		{Name: "Software", Value: strconv.Itoa(image.SoftwareID)},
		{Name: "Image", Value: strconv.Itoa(image.ImageID)},
		{Name: "Default", Value: YesNo(image.IsDefault)},
		{Name: "Position", Value: strconv.Itoa(image.Position)},
		{Name: "Origin", Value: image.Origin},
		{Name: "URL", Value: Text(image.URL)},
		{Name: "Created", Value: Text(image.CreatedAt)},
		{Name: "Updated", Value: Text(image.UpdatedAt)},
	}
}
