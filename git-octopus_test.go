package main

import (
	"os"
)

func ExampleVersionShort() {
	repo := createRepo()
	defer cleanupTestRepo(repo)
	mainWithArgs(repo, "-v")
	// Output: 2.0
}

func initRepo() *repository {
	repo := createRepo()
	repo.git("commit", "--allow-empty", "-m\"first commit\"")
	head := repo.git("rev-parse", "HEAD")
	repo.git("update-ref", "refs/heads/test1", head)
	repo.git("update-ref", "refs/remotes/origin/test1", head)
	repo.git("update-ref", "refs/remotes/origin/test2", head)

	return repo
}

func cleanupTestRepo(repo *repository) error {
	return os.RemoveAll(repo.path)
}

func ExampleOctopusNoPatternGiven() {
	repo := createRepo()
	defer cleanupTestRepo(repo)

	mainWithArgs(repo)
	// Output: Nothing to merge. No pattern given
}

func ExampleOctopusNoBranchMatching() {
	repo := createRepo()
	defer cleanupTestRepo(repo)

	mainWithArgs(repo, "refs/remotes/dumb/*", "refs/remotes/dumber/*")
	// Output: No branch matching "refs/remotes/dumb/* refs/remotes/dumber/*" were found
}
