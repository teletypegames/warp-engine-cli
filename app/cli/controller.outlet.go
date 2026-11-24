package cli

import (
	"errors"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// OutletControllerInterface is the "outlet" commands: the formats the
// catalog leaves the engine in.
type OutletControllerInterface interface {
	// Index lists the outlets.
	Index(cmd *cobra.Command, args []string) error
	// Check shows what one outlet offers and leaves out.
	Check(cmd *cobra.Command, args []string) error
}

// OutletController is the [OutletControllerInterface] over a
// [domain.OutletServiceInterface].
type OutletController struct {
	// OutletService reads the outlets.
	OutletService domain.OutletServiceInterface
}

// Index prints the configured outlets, the catalog API first.
func (c OutletController) Index(cmd *cobra.Command, args []string) error {
	outlets, err := c.OutletService.All()
	if err != nil {
		return err
	}
	return PrintList(cmd, outlets, OutletSerializer{}.ModelToTable(outlets))
}

// Check prints what the outlet KEY offers and leaves out; with --strict it
// fails (exit code 1) when it leaves a software out.
func (c OutletController) Check(cmd *cobra.Command, args []string) error {
	check, err := c.OutletService.Check(args[0])
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, check, OutletSerializer{}.CheckToRecord(*check)); err != nil {
		return err
	}
	if strict, _ := cmd.Flags().GetBool("strict"); strict && len(check.Skipped) > 0 {
		return errors.New("the outlet leaves softwares out")
	}
	return nil
}
