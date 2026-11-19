package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"git.teletypegames.org/engines/warp-engine-cli/app/api"
	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Exit codes of the program; see the package documentation.
const (
	// ExitOK: the command did what it was asked.
	ExitOK = 0
	// ExitError: any failure without a code of its own.
	ExitError = 1
	// ExitUsage: the command line was wrong (unknown flag, missing argument).
	ExitUsage = 2
	// ExitAuth: not signed in, or the token was refused (401) or lacks the
	// scope (403).
	ExitAuth = 3
	// ExitNotFound: no such record, or the resource is off on this server (404).
	ExitNotFound = 4
	// ExitRejected: the server refused the change: a conflict (409) or a
	// record that does not validate (422).
	ExitRejected = 5
	// ExitRateLimited: too many requests (429); wait and try again.
	ExitRateLimited = 6
	// ExitServer: the server failed (5xx) or could not be reached.
	ExitServer = 7
)

// UsageError is a command line that cannot run: a bad flag value, a missing
// argument. It exits with [ExitUsage].
type UsageError struct {
	// Err is what was wrong.
	Err error
}

// Error returns the message of Err.
func (e *UsageError) Error() string {
	return e.Err.Error()
}

// Unwrap returns Err.
func (e *UsageError) Unwrap() error {
	return e.Err
}

// ExitCode returns the exit code for err: [ExitOK] for nil, then by kind —
// [*UsageError], the not-signed-in errors of package domain, the HTTP status
// of an [*api.Error] — and [ExitError] for anything else.
func ExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	var usage *UsageError
	if errors.As(err, &usage) {
		return ExitUsage
	}
	if errors.Is(err, domain.ErrNoProfile) || errors.Is(err, api.ErrNoServer) || errors.Is(err, domain.ErrNoLocalToken) ||
		errors.Is(err, domain.ErrDenied) || errors.Is(err, domain.ErrExpired) {
		return ExitAuth
	}
	status := api.StatusOf(err)
	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return ExitAuth
	case status == http.StatusNotFound:
		return ExitNotFound
	case status == http.StatusConflict || status == http.StatusUnprocessableEntity || status == http.StatusBadRequest:
		return ExitRejected
	case status == http.StatusTooManyRequests:
		return ExitRateLimited
	case status >= 500:
		return ExitServer
	case status == 0 && isNetworkError(err):
		return ExitServer
	}
	return ExitError
}

// isNetworkError reports whether err comes from reaching the server rather
// than from its answer.
func isNetworkError(err error) bool {
	text := err.Error()
	for _, sign := range []string{"connection refused", "no such host", "timeout", "EOF", "connection reset", "tls:"} {
		if strings.Contains(text, sign) {
			return true
		}
	}
	return false
}

// ExactArgs is cobra.ExactArgs whose error is a [*UsageError] that names
// what the arguments are: ExactArgs(1, "SOFTWARE").
func ExactArgs(count int, names ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == count {
			return nil
		}
		if count == 0 {
			return &UsageError{Err: fmt.Errorf("%s takes no arguments", cmd.CommandPath())}
		}
		return &UsageError{Err: fmt.Errorf("%s needs %s (see --help)", cmd.CommandPath(), strings.Join(names, " "))}
	}
}

// RangeArgs is cobra.RangeArgs with a [*UsageError] naming the arguments.
func RangeArgs(minimum int, maximum int, names ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) >= minimum && len(args) <= maximum {
			return nil
		}
		return &UsageError{Err: fmt.Errorf("%s takes %s (see --help)", cmd.CommandPath(), strings.Join(names, " "))}
	}
}

// DeclareListFlags adds the paging flags of a list command — --page,
// --per-page, --all — and --sort with the orders in sorts (no --sort when
// sorts is empty).
func DeclareListFlags(flags *pflag.FlagSet, sorts []string) {
	flags.Int("page", 0, "the page to show, from 1")
	flags.Int("per-page", 0, "records per page (the server caps it)")
	flags.Bool("all", false, "read every page")
	if len(sorts) > 0 {
		flags.String("sort", "", "order: "+strings.Join(sorts, ", ")+"; a leading - reverses it")
	}
}

// FlagsToListQuery reads the flags of [DeclareListFlags].
func FlagsToListQuery(flags *pflag.FlagSet) domain.ListQuery {
	query := domain.ListQuery{}
	query.Page, _ = flags.GetInt("page")
	query.PerPage, _ = flags.GetInt("per-page")
	query.All, _ = flags.GetBool("all")
	query.Sort, _ = flags.GetString("sort")
	return query
}

// FilterFlag is a query parameter of a list command: --platform c64 sends
// platform=c64.
type FilterFlag struct {
	// Flag is the flag name, kebab-case.
	Flag string
	// Key is the query parameter, camelCase.
	Key string
	// Usage is the help text.
	Usage string
	// Bool makes the flag a boolean: --highlighted, --highlighted=false.
	Bool bool
}

// DeclareFilters adds a flag for every filter.
func DeclareFilters(flags *pflag.FlagSet, filters []FilterFlag) {
	for _, filter := range filters {
		if filter.Bool {
			flags.Bool(filter.Flag, false, filter.Usage)
		} else {
			flags.String(filter.Flag, "", filter.Usage)
		}
	}
}

// FlagsToFilter reads the filters that were set into a [domain.Filter].
func FlagsToFilter(flags *pflag.FlagSet, filters []FilterFlag) domain.Filter {
	filter := domain.Filter{}
	for _, item := range filters {
		if !flags.Changed(item.Flag) {
			continue
		}
		if item.Bool {
			value, _ := flags.GetBool(item.Flag)
			filter[item.Key] = YesNoBool(value)
		} else {
			filter[item.Key], _ = flags.GetString(item.Flag)
		}
	}
	return filter
}

// YesNoBool renders value as the API's query boolean, "true" or "false".
func YesNoBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// DeclareYes adds --yes, which answers yes to the question of [Confirm].
func DeclareYes(flags *pflag.FlagSet) {
	flags.BoolP("yes", "y", false, "do not ask for confirmation")
}

// Confirm asks question on the terminal and reports whether the answer was
// yes. With --yes it asks nothing and returns true. Without a terminal on
// standard input it cannot ask: it returns a [*UsageError] telling to pass
// --yes.
func Confirm(cmd *cobra.Command, question string) (bool, error) {
	if yes, _ := cmd.Flags().GetBool("yes"); yes {
		return true, nil
	}
	if !isTerminal(os.Stdin) {
		return false, &UsageError{Err: errors.New("not a terminal, cannot ask: pass --yes to go ahead")}
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "%s [y/N] ", question)
	answer, _ := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes", nil
}

// ErrCancelled is returned when the person answered no to [Confirm].
var ErrCancelled = errors.New("cancelled")

// isTerminal reports whether file is a terminal (a character device).
func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

// readAll reads all of reader; a seam for --data @-.
func readAll(reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}
