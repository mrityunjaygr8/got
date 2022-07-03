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
	Use:   "init",
	Short: "Initialize a new got repository",
	Long:  "Initialize a new got repository",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := got.CreateRepo(args[0])
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("repository: %s created successfully\n", args[0])
	},
}
