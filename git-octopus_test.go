package main

import (
	"io/ioutil"
	"os"
)

func createTestRepo() *repository {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")

	repo := repository{path: dir}

	repo.git("init")
	repo.git("commit", "--allow-empty", "-m\"first commit\"")

	return &repo
}

func cleanupTestRepo(repo *repository) error {
	return os.RemoveAll(repo.path)
}

func ExampleVersionShort() {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

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
