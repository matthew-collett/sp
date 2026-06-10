package shelf

import (
	"fmt"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type shelfDropOpts struct {
	shelf *config.Shelf
	name  string
}

func NewCmdShelfDrop(f *factory.Factory) *cobra.Command {
	opts := &shelfDropOpts{}

	cmd := &cobra.Command{
		Use:               "drop <name>",
		Short:             "Remove an item from the shelf",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cmdutil.ShelfCompletion(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			shelf, err := f.Shelf()
			if err != nil {
				return err
			}
			opts.shelf = shelf
			opts.name = args[0]
			return runCmdShelfDrop(opts)
		},
	}

	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdShelfDrop(opts *shelfDropOpts) error {
	if !opts.shelf.Drop(opts.name) {
		return fmt.Errorf("shelf item not found: %q", opts.name)
	}

	if err := opts.shelf.Write(); err != nil {
		return err
	}

	ui.Success("Dropped %q from shelf", opts.name).Show()
	return nil
}
