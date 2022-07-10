package got

import (
	"errors"
	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set/v2"
)

func Log_graphviz(repo Repo, sha string, seen mapset.Set[string]) {
	if seen.Contains(sha) {
		return
	}
	seen.Add(sha)

	commit_interface, err := Object_read(repo, sha)
	if err != nil {
		log.Fatal(err)
	}

	commit_type := string(commit_interface.get_type())
	if commit_type != "commit" {
		log.Fatal(errors.New("log can only be viewed of a commit"))
	}

	commit, _ := commit_interface.(*commit)
	parents, contains_parent := commit.klvm.Get("parent")

	// base case, this is the initical commit
	if !contains_parent {
		return
	}

	for _, parent := range parents {
		fmt.Printf("c_%s -> c_%s;\n", sha, string(parent))
		Log_graphviz(repo, string(parent), seen)
	}

}
