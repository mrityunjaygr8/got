package got

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

const GIT_DIR = ".git"
const CONFIG_FILE_NAME = "config"
const INITIAL_DESCRIPTION = "Unnamed repository; edit this file 'description' to name the repository.\n"
const INITIAL_HEAD = "ref: refs/heads/master\n"

func isDirectory(path string) (bool, error) {
	exists, err := exists(path)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func (r Repo) repo_path(path ...string) string {
	wrap := []string{r.gitconf_dir}
	wrap = append(wrap, path...)
	return filepath.Join(wrap...)
}

func (r Repo) repo_file(mkdir bool, path ...string) (string, error) {
	_, err := r.repo_dir(mkdir, path[:len(path)-1]...)
	if err != nil {
		return "", err
	}

	return r.repo_path(path...), nil
}

func (r Repo) repo_dir(mkdir bool, path ...string) (string, error) {
	the_path := r.repo_path(path...)

	path_exists, err := exists(the_path)
	if err != nil {
		return "", err
	}

	if path_exists {
		isDir, err := isDirectory(the_path)
		if err != nil {
			return "", err
		}

		if isDir {
			return the_path, nil
		} else {
			return "", fmt.Errorf("Not a directory: %s", the_path)
		}
	}

	if mkdir {
		err := os.MkdirAll(the_path, 0755)
		if err != nil {
			return "", err
		}
		return the_path, nil
	}

	return "", nil

}

func repo_default_config() *viper.Viper {
	myViper := viper.New()
	myViper.SetConfigType("ini")
	myViper.Set("core.repositoryformatversion", 0)
	myViper.Set("core.filemode", false)
	myViper.Set("core.bare", false)

	return myViper
}

func Repo_find(path string, required bool) (Repo, error) {
	abs_path, err := filepath.Abs(path)
	if err != nil {
		return Repo{}, err
	}
	real_path, err := filepath.EvalSymlinks(abs_path)
	if err != nil {
		return Repo{}, err
	}

	here, err := isDirectory(filepath.Join(real_path, GIT_DIR))
	if err != nil {
		return Repo{}, err
	}

	if here {
		return NewRepo(real_path, false)
	}

	parent, err := filepath.EvalSymlinks(filepath.Join(real_path, ".."))
	if err != nil {
		return Repo{}, err
	}

	// Bottom case for recursion If this is true, then we are at the root of the filesystem
	if parent == real_path {
		if required {
			return Repo{}, errors.New("No git repository")
		} else {
			return Repo{}, nil
		}
	}
	return Repo_find(parent, required)
}

func object_resolve(repo Repo, name string) ([]string, error) {
	candidates := make([]string, 0)
	hashRE := regexp.MustCompile("^[0-9A-Fa-f]{4,40}$")
	trimmedName := strings.TrimSpace(name)

	if trimmedName == "" {
		return candidates, errors.New("no name provided")
	}

	if hashRE.Match([]byte(trimmedName)) {
		lowerName := strings.ToLower(trimmedName)
		if len(trimmedName) == 40 {
			candidates = append(candidates, lowerName)
			return candidates, nil
		}

		prefix := name[:2]
		path, err := repo.repo_dir(false, "objects", prefix)
		if err != nil {
			return candidates, err
		}
		rem := name[2:]
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return candidates, err
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), rem) {
				candidates = append(candidates, fmt.Sprintf("%s%s", prefix, file.Name()))
			}
		}
	}

	return candidates, nil
}
