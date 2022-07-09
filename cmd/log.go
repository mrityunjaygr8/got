package cmd

import (
	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:       "log COMMIT-ID",
	Short:     "Display history of a given commit.",
	ValidArgs: []string{"commit-id"},
	Args:      cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commit_id := args[0]
		repo, err := got.NewRepo(".", true)
		if err != nil {
			log.Fatal(err)
		}

		seen := mapset.NewSet[string]()
		fmt.Println("digraph gotlog{")
		got.Log_graphviz(repo, commit_id, seen)
		fmt.Println("}")
	},
}
