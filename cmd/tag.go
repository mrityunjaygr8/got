package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
)

var annotated bool

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.Flags().BoolVarP(&annotated, "annotated", "a", false, "Create an annotated tag")
}

var tagCmd = &cobra.Command{
	Use:   "tag [[-a] NAME [OBJECT]]",
	Short: "List and create tags",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return errors.New("Too many args have been supplied. Only accept a maximum of 2")
		}
		if len(args) == 0 && annotated != false {
			return errors.New("Cannot set annotated flag to be true with no args supplied")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var name, object string
		repo, err := got.NewRepo(".", true)
		if err != nil {
			log.Fatal(err)
		}

		switch len(args) {
		case 0:
		case 1:
			name = args[0]
			object = "HEAD"
		case 2:
			name = args[0]
			object = args[1]
		}

		fmt.Println(repo, name, object)
	},
}
