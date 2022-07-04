package cmd

import (
	"fmt"
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

var object_type string
var write bool

func init() {
	rootCmd.AddCommand(hashObjectCmd)
	hashObjectCmd.Flags().StringVarP(&object_type, "type", "t", "blob", "The type of object")
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "Actually write the object into the database")
}

var hashObjectCmd = &cobra.Command{
	Use:       "hash-object FILE",
	Short:     "Compute object ID and optionally creates a blob from a file",
	ValidArgs: []string{"path"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		var repo got.Repo
		if write {
			my_repo, err := got.NewRepo(".", false)
			if err != nil {
				log.Fatal(err)
			}

			repo = my_repo
		}

		sha, err := got.Object_hash(path, object_type, repo)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%x\n", sha)
		// fmt.Println(sha)
	},
}
