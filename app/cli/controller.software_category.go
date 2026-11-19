package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// SoftwareCategoryControllerInterface is the "category" commands: the vocabulary of a software's category.
type SoftwareCategoryControllerInterface interface {
	// Index lists the software categories.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one category.
	Show(cmd *cobra.Command, args []string) error
	// Create creates a category.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a category.
	Update(cmd *cobra.Command, args []string) error
	// Destroy deletes a category.
	Destroy(cmd *cobra.Command, args []string) error
}

// SoftwareCategoryController is the [SoftwareCategoryControllerInterface] over a
// [domain.SoftwareCategoryServiceInterface].
type SoftwareCategoryController struct {
	// SoftwareCategoryService manages the software categories.
	SoftwareCategoryService domain.SoftwareCategoryServiceInterface
}

// Index prints the software categories matching the filter flags (--needs-review), one
// page or, with --all, every page.
func (c SoftwareCategoryController) Index(cmd *cobra.Command, args []string) error {
	serializer := SoftwareCategorySerializer{}
	page, err := c.SoftwareCategoryService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the category ID.
func (c SoftwareCategoryController) Show(cmd *cobra.Command, args []string) error {
	category_id, err := domain.ParseID("category", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	category, err := c.SoftwareCategoryService.Show(category_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, category, SoftwareCategorySerializer{}.ModelToRecord(*category))
}

// Create creates a category from the field flags (--key and --label are
// needed) and prints it.
func (c SoftwareCategoryController) Create(cmd *cobra.Command, args []string) error {
	serializer := SoftwareCategorySerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	category, err := c.SoftwareCategoryService.Create(attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, category, serializer.ModelToRecord(*category))
}

// Update changes the fields given as flags on the category ID and prints it.
func (c SoftwareCategoryController) Update(cmd *cobra.Command, args []string) error {
	category_id, err := domain.ParseID("category", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	serializer := SoftwareCategorySerializer{}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	category, err := c.SoftwareCategoryService.Update(category_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, category, serializer.ModelToRecord(*category))
}

// Destroy deletes the category ID after asking (not with --yes).
func (c SoftwareCategoryController) Destroy(cmd *cobra.Command, args []string) error {
	category_id, err := domain.ParseID("category", args[0])
	if err != nil {
		return &UsageError{Err: err}
	}
	ok, err := Confirm(cmd, "Delete the category "+args[0]+"?")
	if err != nil || !ok {
		if err == nil {
			err = ErrCancelled
		}
		return err
	}
	if err = c.SoftwareCategoryService.Destroy(category_id); err != nil {
		return err
	}
	PrintDone(cmd, "Category %s deleted.", args[0])
	return nil
}
