package test

import (
	"io/ioutil"
	"lesfurets/git-octopus/git"
	"os"
)

func CreateTempDir() string {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")
	return dir
}

func Cleanup(repo *git.Repository) error {
	return os.RemoveAll(repo.Path)
}