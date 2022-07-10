package cmd

import (
	"fmt"
	"log"

	"github.com/mrtyunjaygr8/got/got"
	"github.com/spf13/cobra"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func init() {
	rootCmd.AddCommand(showRefCmd)
}

var showRefCmd = &cobra.Command{
	Use:   "show-ref",
	Short: "List references",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := got.NewRepo(".", true)
		if err != nil {
			log.Fatal(err)
		}

		ret := orderedmap.New[string, string]()
		refs := got.Ref_list(repo, "", *ret)

		// for pair := refs.Oldest(); pair != nil; pair = pair.Next() {
		// 	fmt.Printf("%s => %s\n", pair.Key, pair.Value)
		// } // prints:

		for _, val := range got.Show_ref(repo, true, "", refs) {
			fmt.Println(val)
		}
	},
}
