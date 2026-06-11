package cmdutil

import (
	"slices"

	"github.com/matthew-collett/sp/internal/factory"
	"github.com/spf13/cobra"
)

const (
	skipConfigCheckAnnotation = "skipConfigCheck"
	skipAuthCheckAnnotation   = "skipAuthCheck"
)

func DisableConfigCheck(cmd *cobra.Command) { disableCheck(cmd, skipConfigCheckAnnotation) }
func DisableAuthCheck(cmd *cobra.Command)   { disableCheck(cmd, skipAuthCheckAnnotation) }

func IsConfigCheckEnabled(cmd *cobra.Command) bool {
	return checkEnabled(cmd, skipConfigCheckAnnotation)
}

func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	return checkEnabled(cmd, skipAuthCheckAnnotation)
}

func disableCheck(cmd *cobra.Command, annotation string) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[annotation] = "true"
}

func checkEnabled(cmd *cobra.Command, annotation string) bool {
	skipCmds := []string{"sp", "help", "completion", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd}
	for c := cmd; c != nil; c = c.Parent() {
		if slices.Contains(skipCmds, c.Name()) {
			return false
		}
		if c.Annotations[annotation] == "true" {
			return false
		}
	}
	return true
}

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
