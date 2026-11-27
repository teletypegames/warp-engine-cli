package cli

import (
	"errors"

	"git.teletypegames.org/engines/warp-engine-cli/app/domain"
	"github.com/spf13/cobra"
)

// FeedControllerInterface is the "feed" commands: the formats the
// catalog leaves the engine in.
type FeedControllerInterface interface {
	// Index lists the feeds.
	Index(cmd *cobra.Command, args []string) error
	// Check shows what one feed offers and leaves out.
	Check(cmd *cobra.Command, args []string) error
}

// FeedController is the [FeedControllerInterface] over a
// [domain.FeedServiceInterface].
type FeedController struct {
	// FeedService reads the feeds.
	FeedService domain.FeedServiceInterface
}

// Index prints the configured feeds, the catalog API first.
func (c FeedController) Index(cmd *cobra.Command, args []string) error {
	feeds, err := c.FeedService.All()
	if err != nil {
		return err
	}
	return PrintList(cmd, feeds, FeedSerializer{}.ModelToTable(feeds))
}

// Check prints what the feed KEY offers and leaves out; with --strict it
// fails (exit code 1) when it leaves a software out.
func (c FeedController) Check(cmd *cobra.Command, args []string) error {
	check, err := c.FeedService.Check(args[0])
	if err != nil {
		return err
	}
	if err = PrintRecord(cmd, check, FeedSerializer{}.CheckToRecord(*check)); err != nil {
		return err
	}
	if strict, _ := cmd.Flags().GetBool("strict"); strict && len(check.Skipped) > 0 {
		return errors.New("the feed leaves softwares out")
	}
	return nil
}
