package cmd

import (
	"fmt"
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(catFileCmd)
}

var catFileCmd = &cobra.Command{
	Use:       "cat-file TYPE OBJECT",
	Short:     "Provide content of repository objects",
	ValidArgs: []string{"type", "object"},
	Args:      cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := got.Repo_find(".", false)
		if err != nil {
			log.Fatal(err)
		}

		obj, err := got.Object_read(repo, got.Object_find(repo, args[1], args[0], true))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(obj.Serialize()))
	},
}
