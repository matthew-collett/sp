package shelf

import (
	"fmt"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/spotify"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type shelfAddOpts struct {
	shelf *config.Shelf
	name  string
	uri   string
}

func NewCmdShelfAdd(f *factory.Factory) *cobra.Command {
	opts := &shelfAddOpts{}

	cmd := &cobra.Command{
		Use:   "add <name> <uri>",
		Short: "Add a Spotify URI to the shelf with a custom name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			shelf, err := f.Shelf()
			if err != nil {
				return err
			}
			opts.shelf = shelf
			opts.name = args[0]
			opts.uri = args[1]
			return runCmdShelfAdd(opts)
		},
	}

	cmdutil.DisableAuthCheck(cmd)

	return cmd
}

func runCmdShelfAdd(opts *shelfAddOpts) error {
	if _, _, err := spotify.ParseURI(opts.uri); err != nil {
		return fmt.Errorf("%q is not a valid Spotify URI", opts.uri)
	}

	if ok := opts.shelf.Add(opts.name, opts.uri); !ok {
		return fmt.Errorf(`%q already exists on shelf. Use "sp shelf drop %s" first`, opts.name, opts.name)
	}

	if err := opts.shelf.Write(); err != nil {
		return err
	}

	ui.Success("Added %q to shelf", opts.name).Show()
	return nil
}
