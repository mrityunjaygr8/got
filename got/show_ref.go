package got

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func ref_resolve(repo Repo, ref string) string {
	var path string
	if strings.HasPrefix(ref, GIT_DIR) {
		path = ref
	} else {
		i_path, err := repo.repo_file(false, ref)
		if err != nil {
			log.Fatal(err)
		}
		path = i_path
	}
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	b := bufio.NewReader(f)
	data, err := b.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	data = data[:len(data)-1]

	if strings.HasPrefix(string(data), "ref: ") {
		return ref_resolve(repo, string(data[5:]))
	} else {
		return string(data)
	}
}

func Ref_list(repo Repo, path string, ret orderedmap.OrderedMap[string, string]) orderedmap.OrderedMap[string, string] {
	if path == "" {
		t_path, err := repo.repo_dir(false, "refs")
		if err != nil {
			log.Fatal(err)
		}

		path = t_path
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		can := filepath.Join(path, file.Name())
		canIsDir, err := isDirectory(can)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(can)
		// strings.TrimPrefix
		if canIsDir {
			ret = Ref_list(repo, can, ret)
		} else {
			ret.Set(strings.TrimPrefix(can, GIT_DIR+"/"), ref_resolve(repo, can))
		}
	}
	return ret
}

func Show_ref(repo Repo, with_hash bool, prefix string, refs orderedmap.OrderedMap[string, string]) []string {
	out := make([]string, 0)
	for pair := refs.Oldest(); pair != nil; pair = pair.Next() {
		var hash string
		if with_hash {
			hash_ := pair.Value + " "
			hash = hash_
		} else {
			hash = ""
		}
		s := fmt.Sprintf("%s%s", hash, pair.Key)
		out = append(out, s)
	}
	return out
}
