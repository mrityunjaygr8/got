package cmd

import (
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:       "init [flags] path",
	Short:     "Initialize a new got repository",
	Long:      "Initialize a new got repository\n" + "Arguments:\n" + "path: The location where the repository is to be initialized",
	ValidArgs: []string{"path"},
	Args:      cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) < 1 {
			path = "."

		} else {
			path = args[0]
		}
		_, err := got.CreateRepo(path)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("repository: %s created successfully\n", path)
	},
}
