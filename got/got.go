package got

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Repo struct {
	worktree_dir string
	gitconf_dir  string
}

func NewRepo(path string, force bool) (Repo, error) {
	repo := &Repo{}
	repo.worktree_dir = path
	repo.gitconf_dir = filepath.Join(path, GIT_DIR)

	isDir, err := isDirectory(repo.gitconf_dir)
	if err != nil {
		return Repo{}, err
	}
	if !(force || isDir) {
		return Repo{}, fmt.Errorf("Not a git repository: %s", repo.gitconf_dir)
	}

	config_file, err := repo.repo_file(false, CONFIG_FILE_NAME)
	if err != nil {
		return Repo{}, err
	}

	config_exists, err := exists(config_file)
	if err != nil {
		return Repo{}, err
	}

	if config_exists {
		viper.SetConfigFile(config_file)
		viper.SetConfigType("ini")
		err = viper.ReadInConfig()
		if err != nil {
			return Repo{}, err
		}
	} else if !force {
		return Repo{}, errors.New("Configuration file missing")
	}

	if !force {
		vers := viper.GetInt("core.repositoryformatversion")
		if vers != 0 {
			return Repo{}, fmt.Errorf("Unsupported repositoryformatversion %d", vers)
		}
	}

	return *repo, nil
}

func CreateRepo(path string) (Repo, error) {
	repo, err := NewRepo(path, true)
	if err != nil {
		return Repo{}, err
	}

	repo_workdir_exists, err := exists(repo.worktree_dir)
	if err != nil {
		return Repo{}, err
	}

	if repo_workdir_exists {
		repo_workdir_is_dir, err := isDirectory(repo.worktree_dir)
		if err != nil {
			return Repo{}, err
		}

		if !repo_workdir_is_dir {
			return Repo{}, fmt.Errorf("%s is not a directory!", repo.worktree_dir)
		}

		repo_workdir_is_empty, err := IsEmpty(repo.worktree_dir)
		if err != nil {
			return Repo{}, err
		}

		if !repo_workdir_is_empty {
			return Repo{}, fmt.Errorf("%s is not empty!", repo.worktree_dir)
		}
	} else {
		os.Mkdir(repo.worktree_dir, 0755)
	}

	_, err = repo.repo_dir(true, "branches")
	if err != nil {
		return Repo{}, err
	}
	_, err = repo.repo_dir(true, "objects")
	if err != nil {
		return Repo{}, err
	}
	_, err = repo.repo_dir(true, "refs", "tags")
	if err != nil {
		return Repo{}, err
	}
	_, err = repo.repo_dir(true, "refs", "heads")
	if err != nil {
		return Repo{}, err
	}

	description, err := repo.repo_file(true, "description")
	if err != nil {
		return Repo{}, err
	}
	err = os.WriteFile(description, []byte(INITIAL_DESCRIPTION), 0644)
	if err != nil {
		return Repo{}, err
	}

	head, err := repo.repo_file(true, "HEAD")
	if err != nil {
		return Repo{}, err
	}
	err = os.WriteFile(head, []byte(INITIAL_HEAD), 0644)
	if err != nil {
		return Repo{}, err
	}

	config, err := repo.repo_file(true, CONFIG_FILE_NAME)
	if err != nil {
		return Repo{}, err
	}
	viper := repo_default_config()
	viper.SafeWriteConfigAs(config)

	return repo, nil
}
