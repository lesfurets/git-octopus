package test

import (
	"github.com/lesfurets/git-octopus/git"
	"io/ioutil"
	"os"
)

func CreateTempDir() string {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")
	return dir
}

func Cleanup(repo *git.Repository) error {
	return os.RemoveAll(repo.Path)
}
