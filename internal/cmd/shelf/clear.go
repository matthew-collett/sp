package shelf

import (
	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type shelfClearOpts struct {
	shelf *config.Shelf
}

func NewCmdShelfClear(f *factory.Factory) *cobra.Command {
	opts := &shelfClearOpts{}

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Remove all items from the shelf",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			shelf, err := f.Shelf()
			if err != nil {
				return err
			}
			opts.shelf = shelf
			return runCmdShelfClear(opts)
		},
	}

	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdShelfClear(opts *shelfClearOpts) error {
	if opts.shelf.IsEmpty() {
		ui.Info("Shelf is already empty").Show()
		return nil
	}

	n := len(opts.shelf.Items)
	ok, err := ui.PromptConfirm("Remove all %d shelf items?", n)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	opts.shelf.Clear()
	if err := opts.shelf.Write(); err != nil {
		return err
	}

	ui.Success("Cleared %d shelf items", n).Show()
	return nil
}
