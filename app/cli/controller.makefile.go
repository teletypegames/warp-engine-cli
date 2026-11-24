package cli

import (
	"fmt"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// MakefileControllerInterface is "check makefile".
type MakefileControllerInterface interface {
	// Check checks a project's Makefile against the contract.
	Check(cmd *cobra.Command, args []string) error
}

// MakefileController is the [MakefileControllerInterface] over a
// [domain.MakefileServiceInterface].
type MakefileController struct {
	// MakefileService checks the Makefile.
	MakefileService domain.MakefileServiceInterface
	// AuthService connects, for the server's platforms.
	AuthService domain.AuthServiceInterface
}

// Check checks the Makefile in DIR (the current directory without it),
// prints every required target with its state, and fails when one is
// missing or fails. It reads the platforms from the server [ConnectBuild]
// finds; without one it checks the targets every project has and says so.
func (c MakefileController) Check(cmd *cobra.Command, args []string) error {
	directory := "."
	if len(args) == 1 {
		directory = args[0]
	}
	_ = ConnectBuild(cmd, c.AuthService)
	platform, _ := cmd.Flags().GetString("platform")
	app, _ := cmd.Flags().GetString("app")
	report, err := c.MakefileService.Check(directory, platform, app)
	if err != nil {
		return err
	}
	serializer := MakefileSerializer{}
	if outputOf(cmd) == OutputTable {
		fmt.Fprintln(cmd.ErrOrStderr(), serializer.Heading(*report))
		if report.Note != "" {
			fmt.Fprintln(cmd.ErrOrStderr(), "Note: "+report.Note)
		}
	}
	if err = PrintList(cmd, report, serializer.ReportToTable(*report)); err != nil {
		return err
	}
	if broken := report.Broken(); len(broken) > 0 {
		return fmt.Errorf("the Makefile breaks the contract: %s", strings.Join(broken, ", "))
	}
	return nil
}
