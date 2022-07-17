package got

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type object interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte)
	get_type() []byte
	get_repo() Repo
}

func Object_write(o object, actually_write bool) []byte {
	data, err := o.Serialize()
	if err != nil {
		log.Fatal(err)
	}

	result := make([]byte, 0)
	result = append(result, o.get_type()...)
	result = append(result, ' ')
	result = append(result, []byte(fmt.Sprint(len(data)))...)
	result = append(result, '\x00')
	result = append(result, data...)

	h := sha1.New()
	fmt.Fprint(h, result)
	sha := h.Sum(nil)

	if actually_write {
		repo := o.get_repo()
		sha_string := fmt.Sprintf("%x", sha)
		path, err := repo.repo_file(actually_write, "objects", sha_string[:2], sha_string[2:])
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		w := zlib.NewWriter(f)
		w.Write(result)
		w.Close()
	}
	return sha
}

func Object_read(repo Repo, sha string) (object, error) {
	path, err := repo.repo_file(false, "objects", sha[:2], sha[2:])
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := zlib.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	raw_buf := bytes.NewBuffer(make([]byte, 0))
	_, err = io.Copy(raw_buf, r)
	if err != nil {
		return nil, err
	}

	raw := raw_buf.Bytes()

	x := bytes.Index(raw, []byte(" "))
	format := string(raw[0:x])

	y := bytes.Index(raw, []byte("\x00"))
	size_str := string(raw[x+1 : y])
	size, err := strconv.Atoi(size_str)
	if err != nil {
		return nil, err
	}

	if size != len(raw)-y-1 {
		return nil, fmt.Errorf("Malformed object %s: bad length", sha)
	}

	var c object
	switch format {
	case "blob":
		c = blob{repo: repo, blobdata: raw[y+1:]}
	case "commit":
		c = &commit{repo: repo}
		c.Deserialize(raw[y+1:])
	case "tree":
		c = &tree{repo: repo}
		c.Deserialize(raw[y+1:])
	default:
		return nil, fmt.Errorf("Unknown type %s for object %s", format, sha)

	}

	return c, nil
}

func Object_find(repo Repo, name string, format string, follow bool) (string, error) {
	sha, err := object_resolve(repo, name)
	if err != nil {
		return "", err
	}

	if len(sha) > 1 {
		var many strings.Builder
		many.WriteString(fmt.Sprintf("Ambiguous reference %s: Candidates are:\n", name))
		for _, name := range sha {
			many.WriteString(fmt.Sprintf(" - %s\n", name))
		}

		return many.String(), nil
	}

	newSha := sha[0]

	if format == "" {
		return newSha, nil
	}

	for {
		obj, err := Object_read(repo, newSha)
		if err != nil {
			return "", err
		}

		if string(obj.get_type()) == format {
			return newSha, nil
		}

		if !follow {
			return "", nil
		}

		if string(obj.get_type()) == "tag" {
			theTag, _ := obj.(*tag)
			newShaList, _ := theTag.klvm.Get("object")
			newSha = string(newShaList[0])
		} else if string(obj.get_type()) == "commit" && format == "tree" {
			theCommit, _ := obj.(*commit)
			newShaList, _ := theCommit.klvm.Get("tree")
			newSha = string(newShaList[0])
		} else {
			return "", nil
		}
	}
}

func Object_hash(path string, format string, repo Repo) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make([]byte, 0)
	tmp := make([]byte, 16)
	reader := bufio.NewReader(f)
	for {
		n, err := reader.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		data = append(data, tmp[0:n]...)
	}

	var obj object
	switch format {
	case "blob":
		obj = blob{repo: repo, blobdata: data}
	case "commit":
		obj := &commit{repo: repo}
		obj.Deserialize(data)
	case "tree":
		obj := &tree{repo: repo}
		obj.Deserialize(data)
	default:
		return nil, fmt.Errorf("Unknown type %s", format)
	}

	return Object_write(obj, false), nil

}
