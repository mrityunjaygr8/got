package cmd

import (
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsTreeCmd)
}

var lsTreeCmd = &cobra.Command{
	Use:       "ls-tree OBJECT-ID",
	Short:     "Pretty-print a tree object.",
	ValidArgs: []string{"object_id"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		object_id := args[0]
		repo, err := got.NewRepo(".", true)
		if err != nil {
			log.Fatal(err)
		}

		obj, err := got.Object_read(repo, got.Object_find(repo, object_id, "tree", false))
		if err != nil {
			log.Fatal(err)
		}

		got.Ls_tree(obj)
	},
}
