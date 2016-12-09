package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
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

func TestOctopusCommitConfigError(t *testing.T) {
	repo := createTestRepo()

	repo.git("config", "octopus.commit", "bad_value")

	err := mainWithArgs(repo, "-v")

	assert.NotNil(t, err)
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

func TestOctopusAlreadyUpToDate(t *testing.T) {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	repo.writeFile("foo", "First line")
	repo.git("add", "foo")
	repo.git("commit", "-m\"first commit\"")
	// Create a branch on this first commit
	repo.git("branch", "outdated_branch")
	expected, _ := repo.git("rev-parse", "HEAD")

	mainWithArgs(repo, "outdated_branch")

	actual, _ := repo.git("rev-parse", "HEAD")

	assert.Equal(t, expected, actual)
}
