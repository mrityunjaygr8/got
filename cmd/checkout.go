package cmd

import (
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

var checkoutCmd = &cobra.Command{
	Use:       "checkout COMMIT PATH",
	Short:     "Checkout a commit inside of a directory",
	ValidArgs: []string{"commit", "path"},
	Args:      cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		commit_id := args[0]
		path := args[1]
		repo, err := got.Repo_find(".", false)
		if err != nil {
			log.Fatal(err)
		}

		o, err := got.Object_find(repo, commit_id, "commit", true)
		if err != nil {
			log.Fatal(err)
		}
		obj, err := got.Object_read(repo, o)
		if err != nil {
			log.Fatal(err)
		}

		got.Tree_checkout(repo, obj, path)
	},
}
