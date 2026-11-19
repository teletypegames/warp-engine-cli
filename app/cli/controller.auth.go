package cli

import (
	"errors"
	"fmt"
	"os"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// Session is what the root command learned before the command ran: the
// profile it connected with. It is shared by pointer between the root and
// the controllers that need it.
type Session struct {
	// Profile is the connected profile, nil before connecting or for a
	// command that does not connect.
	Profile *domain.ProfileModel
}

// AuthControllerInterface is the sign-in commands: login, logout, whoami,
// server, and the connection every other command makes first.
type AuthControllerInterface interface {
	// Connect picks the server and token from the global flags and the
	// environment; the root command runs it before every online command.
	Connect(cmd *cobra.Command, args []string) error
	// Login signs in to a server and saves the profile.
	Login(cmd *cobra.Command, args []string) error
	// Logout forgets a profile, revoking its token.
	Logout(cmd *cobra.Command, args []string) error
	// Whoami shows who the server takes the token for.
	Whoami(cmd *cobra.Command, args []string) error
	// Server shows what the connected engine says about itself.
	Server(cmd *cobra.Command, args []string) error
}

// AuthController is the [AuthControllerInterface] over an
// [domain.AuthServiceInterface].
type AuthController struct {
	// AuthService signs in and connects.
	AuthService domain.AuthServiceInterface
	// Session receives the connected profile.
	Session *Session
	// OpenBrowser opens a URL in the person's browser; nil opens none.
	OpenBrowser func(url string) error
}

// Connect reads --profile, --server and --token (or WARP_ENGINE_PROFILE,
// WARP_ENGINE_URL, WARP_ENGINE_TOKEN; the flags win) and connects through
// [domain.AuthService.Connect], storing the profile in the session.
func (c AuthController) Connect(cmd *cobra.Command, args []string) error {
	selection := domain.Selection{
		Profile: flagOrEnv(cmd, "profile", "WARP_ENGINE_PROFILE"),
		Server:  flagOrEnv(cmd, "server", "WARP_ENGINE_URL"),
		Token:   flagOrEnv(cmd, "token", "WARP_ENGINE_TOKEN"),
	}
	profile, err := c.AuthService.Connect(selection)
	if err != nil {
		return err
	}
	c.Session.Profile = profile
	return nil
}

// Login signs in to URL (WarpEngine Local, [domain.LocalURL], without one)
// and saves the profile as the current one. The way depends on the flags and
// on what the server offers:
//
//   - --token TOKEN: the token is checked (GET /me) and saved;
//   - --local, or a server that offers the token file (WarpEngine Local):
//     a profile that reads the local token file on every run;
//   - otherwise a device sign-in: the code and the approval page are
//     printed, the page is opened in the browser (not with --no-browser),
//     and the command waits until it is approved there.
//
// --name names the profile (default: the host, or "local"); --scopes asks a
// device sign-in for other scopes than admin:read and admin:write.
func (c AuthController) Login(cmd *cobra.Command, args []string) error {
	server_url := domain.LocalURL
	if len(args) == 1 {
		server_url = args[0]
	}
	name, _ := cmd.Flags().GetString("name")
	token := flagOrEnv(cmd, "token", "WARP_ENGINE_TOKEN")
	local, _ := cmd.Flags().GetBool("local")
	scopes, _ := cmd.Flags().GetStringSlice("scopes")
	no_browser, _ := cmd.Flags().GetBool("no-browser")

	var profile *domain.ProfileModel
	var err error
	switch {
	case token != "":
		profile, err = c.AuthService.FinishLogin(server_url, token, name)
	case local:
		profile, err = c.AuthService.LoginLocal(server_url, name)
	default:
		var server *domain.ServerModel
		server, err = c.AuthService.Discover(server_url)
		if err != nil {
			return fmt.Errorf("cannot reach %s: %w", server_url, err)
		}
		if server.Supports(domain.SignInTokenFile) {
			profile, err = c.AuthService.LoginLocal(server_url, name)
		} else {
			profile, err = c.deviceLogin(cmd, server_url, name, scopes, no_browser)
		}
	}
	if err != nil {
		return err
	}
	PrintDone(cmd, "Signed in to %s as %s; profile %q is now current.", profile.URL, labelOf(profile), profile.Name)
	return nil
}

// deviceLogin runs a device sign-in at server_url and saves its token.
func (c AuthController) deviceLogin(cmd *cobra.Command, server_url string, name string, scopes []string, no_browser bool) (*domain.ProfileModel, error) {
	host, _ := os.Hostname()
	grant, err := c.AuthService.StartDeviceLogin(server_url, "warp-engine-cli on "+host, scopes)
	if err != nil {
		return nil, err
	}
	out := cmd.ErrOrStderr()
	fmt.Fprintf(out, "To sign in, open\n\n    %s\n\nand check that it shows the code %s.\n", grant.VerificationURLComplete, grant.UserCode)
	if !no_browser && c.OpenBrowser != nil {
		if c.OpenBrowser(grant.VerificationURLComplete) == nil {
			fmt.Fprintln(out, "(The page was opened in your browser.)")
		}
	}
	fmt.Fprintln(out, "Waiting for the approval …")
	token, err := c.AuthService.WaitForApproval(*grant)
	if err != nil {
		return nil, err
	}
	return c.AuthService.FinishLogin(server_url, token, name)
}

// Logout forgets the profile PROFILE (the current one without it). Its
// token is revoked on the server first, unless --keep-token; a Local
// profile has no token of its own and is only forgotten.
func (c AuthController) Logout(cmd *cobra.Command, args []string) error {
	name := ""
	if len(args) == 1 {
		name = args[0]
	}
	if name == "" {
		name = flagOrEnv(cmd, "profile", "WARP_ENGINE_PROFILE")
	}
	keep, _ := cmd.Flags().GetBool("keep-token")
	profile, err := c.AuthService.Logout(name, !keep)
	if err != nil {
		if errors.Is(err, domain.ErrNoProfile) && name != "" {
			return fmt.Errorf("no profile named %q (see \"warp-engine profile list\")", name)
		}
		return err
	}
	PrintDone(cmd, "Signed out of %s; profile %q removed.", profile.URL, profile.Name)
	return nil
}

// Whoami prints the principal of the connected token: label, owner,
// scopes, expiry and the resources it may use.
func (c AuthController) Whoami(cmd *cobra.Command, args []string) error {
	principal, err := c.AuthService.Whoami()
	if err != nil {
		return err
	}
	return PrintRecord(cmd, principal, PrincipalSerializer{}.ModelToRecord(*principal, c.Session.Profile))
}

// Server prints what the connected engine says at /.well-known/warp-engine:
// name, version, profile, sign-in methods, Admin API path.
func (c AuthController) Server(cmd *cobra.Command, args []string) error {
	server_url := c.Session.Profile.URL
	server, err := c.AuthService.Discover(server_url)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, server, ServerSerializer{}.ModelToRecord(*server, server_url))
}

// labelOf names who a profile signs in as: its label, or "the local token".
func labelOf(profile *domain.ProfileModel) string {
	if profile.Label != "" {
		return profile.Label
	}
	if profile.TokenFile {
		return "the local token"
	}
	return "the token"
}

// flagOrEnv returns the flag name of cmd when it was set, otherwise the
// environment variable env.
func flagOrEnv(cmd *cobra.Command, name string, env string) string {
	if flag := cmd.Flags().Lookup(name); flag != nil && flag.Changed {
		return flag.Value.String()
	}
	return os.Getenv(env)
}
