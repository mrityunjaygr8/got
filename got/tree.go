package got

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"strings"
)

type treeLeaf struct {
	mode []byte
	path []byte
	sha  []byte
}

type tree struct {
	repo  Repo
	leafs []treeLeaf
}

func (t tree) Serialize() ([]byte, error) {
	return tree_serialize(t.leafs)
}

func (t *tree) Deserialize(data []byte) {
	t.leafs = tree_parse(data)
}

func (t tree) get_type() []byte {
	return []byte("tree")
}

func (t tree) get_repo() Repo {
	return t.repo
}

func NewTreeLeaf(mode, path, sha []byte) *treeLeaf {
	return &treeLeaf{mode: mode, path: path, sha: sha}
}

func tree_parse_one(raw []byte, start int) (int, *treeLeaf) {
	x := bytes.IndexByte(raw[start:], ' ') + start
	if x-start != 5 && x-start != 6 {
		log.Fatal(errors.New("malformed tree leaf"))
	}

	mode := raw[start:x]

	y := bytes.IndexByte(raw[x:], 0) + x
	path := raw[x+1 : y]

	sha := []byte(fmt.Sprintf("%x", raw[y+1:y+21]))

	return y + 21, NewTreeLeaf(mode, path, sha)
}

func tree_parse(raw []byte) []treeLeaf {

	pos := 0
	max := binary.Size(raw)

	ret := make([]treeLeaf, 0)
	for pos < max {
		in, data := tree_parse_one(raw, pos)
		ret = append(ret, *data)
		pos = in
	}

	return ret
}

func tree_serialize(leaves []treeLeaf) ([]byte, error) {
	raw := make([]byte, 0)

	for _, leaf := range leaves {
		raw = append(raw, leaf.mode...)
		raw = append(raw, ' ')
		raw = append(raw, leaf.path...)
		raw = append(raw, 0)
		sha := fmt.Sprintf("%x", leaf.sha)
		raw = append(raw, []byte(sha)...)
	}

	return raw, nil
}

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
