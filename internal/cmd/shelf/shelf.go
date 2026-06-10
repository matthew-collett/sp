package shelf

import (
	"sort"

	"github.com/matthew-collett/sp/internal/cmdutil"
	"github.com/matthew-collett/sp/internal/config"
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/matthew-collett/sp/internal/ui"
	"github.com/spf13/cobra"
)

type shelfOpts struct {
	shelf   *config.Shelf
	noPager bool
}

func NewCmdShelf(f *factory.Factory) *cobra.Command {
	opts := &shelfOpts{}

	cmd := &cobra.Command{
		Use:     "shelf",
		Aliases: []string{"sh"},
		Short:   "List saved quick-play shortcuts",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			shelf, err := f.Shelf()
			if err != nil {
				return err
			}
			opts.shelf = shelf
			return runCmdShelf(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.noPager, "no-pager", false, "disable pager for output")

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(
		NewCmdShelfAdd(f),
		NewCmdShelfDrop(f),
		NewCmdShelfRename(f),
	)

	return cmd
}

func runCmdShelf(opts *shelfOpts) error {
	if opts.shelf.IsEmpty() {
		ui.Info(`Shelf is empty. Use "sp shelf add <name> <uri>" to add items.`).Show()
		return nil
	}

	t := ui.NewTable()
	t.Title(ui.Text("shelf items"))
	t.Header("NAME", "TYPE", "URI")

	names := make([]string, 0, len(opts.shelf.Items))
	for name := range opts.shelf.Items {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		item := opts.shelf.Items[name]
		t.Row(
			ui.Text(name).Bold(),
			ui.Text(item.Type).Dimmed().Fixed(),
			ui.Text(item.URI).Yellow().Fixed(),
		)
	}

	t.Render(!opts.noPager)
	return nil
}
