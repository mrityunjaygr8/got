package got

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"strings"
)

func Ls_tree(obj object) {
	if string(obj.get_type()) != "tree" {
		log.Fatal(errors.New("Not a tree object"))
	}
	theTree, _ := obj.(*tree)
	for _, leaf := range theTree.leafs {
		mode := strings.Repeat("0", 6-binary.Size(leaf.mode)) + string(leaf.mode)
		path, err := Object_read(theTree.repo, string(leaf.sha))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s %s\t%s\n", mode, string(path.get_type()), leaf.sha, leaf.path)
	}
}
