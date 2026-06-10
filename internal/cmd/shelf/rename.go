package shelf

import (
	"fmt"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type shelfRenameOpts struct {
	shelf *config.Shelf
	old   string
	new   string
}

func NewCmdShelfRename(f *factory.Factory) *cobra.Command {
	opts := &shelfRenameOpts{}

	cmd := &cobra.Command{
		Use:               "rename <old> <new>",
		Short:             "Rename a shelf item",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: cmdutil.ShelfCompletion(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			shelf, err := f.Shelf()
			if err != nil {
				return err
			}
			opts.shelf = shelf
			opts.old = args[0]
			opts.new = args[1]
			return runCmdShelfRename(opts)
		},
	}

	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdShelfRename(opts *shelfRenameOpts) error {
	if !opts.shelf.Has(opts.old) {
		return fmt.Errorf("shelf item not found: %q", opts.old)
	}

	if opts.shelf.Has(opts.new) {
		return fmt.Errorf("%q already exists on shelf", opts.new)
	}

	if !opts.shelf.Rename(opts.old, opts.new) {
		return fmt.Errorf("failed to rename %q", opts.old)
	}

	if err := opts.shelf.Write(); err != nil {
		return err
	}

	ui.Success("Renamed %q to %q", opts.old, opts.new).Show()
	return nil
}
