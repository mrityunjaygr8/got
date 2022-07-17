package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

var got_type string

func init() {
	rootCmd.AddCommand(revParseCmd)
	revParseCmd.Flags().StringVarP(&got_type, "type", "t", "", "Specify the expected type")
}

var revParseCmd = &cobra.Command{
	Use:       "rev-parse NAME",
	Short:     "Parse revision (or other objects )identifiers",
	ValidArgs: []string{"name"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("only one argument is expected")
		}

		if got_type != "" {
			if got_type == "blob" || got_type == "commit" || got_type == "tag" || got_type == "tree" {
			} else {
				return errors.New("invalid value of got_type. Valid values are blob, commit, tree and tag")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		repo, err := got.NewRepo(".", true)
		if err != nil {
			log.Fatal(err)
		}

		obj, err := got.Object_find(repo, name, got_type, true)
		fmt.Println(obj)
	},
}
