package cli

import (
	"os/exec"
	"runtime"
)

// OpenBrowser opens address in the default browser of the desktop: open on
// macOS, rundll32 on Windows, xdg-open elsewhere. It returns as soon as the
// opener started; an error means there is no opener (a server, a container).
func OpenBrowser(address string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", address)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", address)
	default:
		command = exec.Command("xdg-open", address)
	}
	return command.Start()
}
