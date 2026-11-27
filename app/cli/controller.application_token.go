package cli

import (
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// ApplicationTokenControllerInterface is the "token" commands: the tokens
// programs sign in to the engine with.
type ApplicationTokenControllerInterface interface {
	// Index lists the tokens.
	Index(cmd *cobra.Command, args []string) error
	// Show shows one token.
	Show(cmd *cobra.Command, args []string) error
	// Create issues a token and prints it once.
	Create(cmd *cobra.Command, args []string) error
	// Update changes a token.
	Update(cmd *cobra.Command, args []string) error
	// Destroy revokes a token.
	Destroy(cmd *cobra.Command, args []string) error
	// Rotate replaces a token and prints the new one once.
	Rotate(cmd *cobra.Command, args []string) error
}

// ApplicationTokenController is the [ApplicationTokenControllerInterface]
// over a [domain.ApplicationTokenServiceInterface].
type ApplicationTokenController struct {
	// ApplicationTokenService manages the tokens.
	ApplicationTokenService domain.ApplicationTokenServiceInterface
}

// applicationTokenHint tells what to do with a new application token.
const applicationTokenHint = "Sign in with it: warp-engine login <server> --token <token>, or set WARP_ENGINE_TOKEN."

// Index prints the tokens matching --state and --kind, one page or, with
// --all, every page.
func (c ApplicationTokenController) Index(cmd *cobra.Command, args []string) error {
	serializer := ApplicationTokenSerializer{}
	page, err := c.ApplicationTokenService.All(serializer.FlagsToFilter(cmd.Flags()), FlagsToListQuery(cmd.Flags()))
	if err != nil {
		return err
	}
	if err = PrintList(cmd, page.Items, serializer.ModelToTable(page.Items)); err != nil {
		return err
	}
	PrintPaging(cmd, page.Page, page.TotalPages, page.Total)
	return nil
}

// Show prints the token ID (never its plain token, which the server no
// longer knows).
func (c ApplicationTokenController) Show(cmd *cobra.Command, args []string) error {
	application_token_id, err := argID("token", args[0])
	if err != nil {
		return err
	}
	token, err := c.ApplicationTokenService.Show(application_token_id)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, token, ApplicationTokenSerializer{}.ModelToRecord(*token))
}

// Create issues a token from --name, --scopes (both needed), --unrestricted,
// --expires-at and --owner-id, and prints it with its plain token, the
// only time that is shown.
func (c ApplicationTokenController) Create(cmd *cobra.Command, args []string) error {
	serializer := ApplicationTokenSerializer{}
	attributes, err := serializer.CreateToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	token, err := c.ApplicationTokenService.Create(attributes)
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, token, serializer.ModelToRecord(*token)); err != nil {
		return err
	}
	printPlainToken(cmd, token.PlainToken, applicationTokenHint)
	return nil
}

// Update changes the fields given as flags on the token ID and prints it.
func (c ApplicationTokenController) Update(cmd *cobra.Command, args []string) error {
	serializer := ApplicationTokenSerializer{}
	application_token_id, err := argID("token", args[0])
	if err != nil {
		return err
	}
	attributes, err := serializer.RequestToAttributes(cmd.Flags())
	if err != nil {
		return err
	}
	token, err := c.ApplicationTokenService.Update(application_token_id, attributes)
	if err != nil {
		return err
	}
	return PrintRecord(cmd, token, serializer.ModelToRecord(*token))
}

// Destroy revokes the token ID after asking (not with --yes).
func (c ApplicationTokenController) Destroy(cmd *cobra.Command, args []string) error {
	application_token_id, err := argID("token", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Revoke the token "+args[0]+"? Whatever signs in with it stops working."); err != nil {
		return err
	}
	if err = c.ApplicationTokenService.Destroy(application_token_id); err != nil {
		return err
	}
	PrintDone(cmd, "Token %d revoked.", application_token_id)
	return nil
}

// Rotate replaces the token ID after asking (not with --yes): the old one
// stops working, the new one is pushed to its CI repositories and printed once.
func (c ApplicationTokenController) Rotate(cmd *cobra.Command, args []string) error {
	application_token_id, err := argID("token", args[0])
	if err != nil {
		return err
	}
	if err = confirmOrCancel(cmd, "Rotate the token "+args[0]+"? The current one stops working."); err != nil {
		return err
	}
	token, err := c.ApplicationTokenService.Rotate(application_token_id)
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, token, ApplicationTokenSerializer{}.ModelToRecord(*token)); err != nil {
		return err
	}
	printPlainToken(cmd, token.PlainToken, applicationTokenHint)
	return nil
}
