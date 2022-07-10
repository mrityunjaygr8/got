package got

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func tree_checkout(repo Repo, thisTree tree, path string) {
	for _, leaf := range thisTree.leafs {
		obj, err := Object_read(repo, string(leaf.sha))
		if err != nil {
			log.Fatal(err)
		}

		dest := filepath.Join(path, string(leaf.path))
		objType := string(obj.get_type())

		if objType == "tree" {
			os.Mkdir(dest, 0755)
			treeObj := obj.(*tree)
			tree_checkout(repo, *treeObj, dest)
		} else if objType == "blob" {
			f, err := os.Create(dest)
			defer f.Close()
			if err != nil {
				log.Fatal(err)
			}

			blobObj := obj.(blob)
			fmt.Fprint(f, string(blobObj.blobdata))
		}
	}
}

func Tree_checkout(repo Repo, obj object, path string) {
	if string(obj.get_type()) == "commit" {
		commitObj := obj.(*commit)
		tree, ok := commitObj.klvm.Get("tree")
		if !ok {
			log.Fatal(errors.New("tree not found in commit details"))
		}
		theObj, err := Object_read(repo, string(tree[0]))
		if err != nil {
			log.Fatal(err)
		}

		obj = theObj
	}
	finalTree := *obj.(*tree)

	exists, err := exists(path)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		isDir, err := isDirectory(path)
		if err != nil {
			log.Fatal(err)
		}

		if !isDir {
			log.Fatal(fmt.Errorf("Not a directory %s!", path))
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}

		if len(files) > 0 {
			log.Fatal(fmt.Errorf("Not empty %s!", path))
		}
	} else {
		os.MkdirAll(path, 0755)
	}

	tree_checkout(repo, finalTree, path)
}
