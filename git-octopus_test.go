package main

import (
	"io/ioutil"
	"os"
	"github.com/lesfurets/git-octopus/git"
)

func createTestRepo() *git.Repository {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")

	repo := git.Repository{Path: dir}

	repo.Git("init")
	repo.Git("commit", "--allow-empty", "-m\"first commit\"")

	return &repo
}

func cleanupTestRepo(repo *git.Repository) error {
	return os.RemoveAll(repo.Path)
}

func ExampleVersionShort() {
	repo := createTestRepo()
	mainWithArgs(repo, "-v")
	// Output: 2.0
}

func ExampleOctopusNoPatternGiven() {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	mainWithArgs(repo)
	// Output: Nothing to merge. No pattern given
}

func ExampleOctopusNoBranchMatching() {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	mainWithArgs(repo, "refs/remotes/dumb/*", "refs/remotes/dumber/*")
	// Output: No branch matching "refs/remotes/dumb/* refs/remotes/dumber/*" were found
}
