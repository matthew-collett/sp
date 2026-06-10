package cmdutil

import (
	"github.com/matthew-collett/sp/internal/factory"
	"github.com/spf13/cobra"
)

func AddGroup(parent *cobra.Command, title string, cmds ...*cobra.Command) {
	g := &cobra.Group{
		Title: title,
		ID:    title,
	}
	parent.AddGroup(g)
	for _, c := range cmds {
		c.GroupID = g.ID
		parent.AddCommand(c)
	}
}

func ShelfCompletion(f *factory.Factory) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		shelf, err := f.Shelf()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := make([]string, 0, len(shelf.Items))
		for name := range shelf.Items {
			names = append(names, name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	}
}
