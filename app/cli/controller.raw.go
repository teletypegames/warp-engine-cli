package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// RawControllerInterface is the "api" command: any request, answer as it is.
type RawControllerInterface interface {
	// Send sends METHOD PATH and prints the answer.
	Send(cmd *cobra.Command, args []string) error
}

// RawController is the [RawControllerInterface] over a
// [domain.RawServiceInterface].
type RawController struct {
	// RawService sends the request.
	RawService domain.RawServiceInterface
}

// Send sends METHOD PATH with the --query parameters (key=value, repeatable)
// and the --data body (JSON, or @file.json), and prints the JSON answer
// indented, whatever --output says. An empty answer (204) prints nothing.
func (c RawController) Send(cmd *cobra.Command, args []string) error {
	pairs, _ := cmd.Flags().GetStringArray("query")
	query := url.Values{}
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || key == "" {
			return &UsageError{Err: fmt.Errorf("--query takes key=value, not %q", pair)}
		}
		query.Add(key, value)
	}
	var body any
	if data, _ := cmd.Flags().GetString(DataFlag); data != "" {
		parsed, err := readData(data)
		if err != nil {
			return err
		}
		body = parsed
	}
	answer, err := c.RawService.Send(args[0], args[1], query, body)
	if err != nil {
		return err
	}
	if len(answer) == 0 {
		return nil
	}
	var indented bytes.Buffer
	if err = json.Indent(&indented, answer, "", "  "); err != nil {
		indented.Reset()
		indented.Write(answer)
	}
	fmt.Fprintln(cmd.OutOrStdout(), indented.String())
	return nil
}
