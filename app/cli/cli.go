package cli

import (
	"fmt"
	"io"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
)

// Start builds the command tree around the services of d, runs it with
// args (without the program name) and returns the exit code. Errors are
// printed to stderr as "warp-engine: <message>".
func Start(d domain.Domain, args []string, stdout io.Writer, stderr io.Writer) int {
	session := &Session{}
	root := Router{
		AuthController:     AuthController{AuthService: d.AuthService, Session: session, OpenBrowser: OpenBrowser},
		ProfileController:  ProfileController{ProfileService: d.ProfileService},
		SoftwareController: SoftwareController{SoftwareService: d.SoftwareService},
		RawController:      RawController{RawService: d.RawService},
	}.Init()
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	err := root.Execute()
	if err != nil {
		fmt.Fprintln(stderr, "warp-engine:", err)
	}
	return ExitCode(err)
}
