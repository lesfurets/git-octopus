package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func (repo *repository) writeFile(name string, lines ...string) {
	fileName := filepath.Join(repo.path, name)
	ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

func ExampleVersionShort() {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	mainWithArgs(repo, "-v")
	// Output: 2.0
}

func TestOctopusCommitConfigError(t *testing.T) {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

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

func TestOctopus3branches(t *testing.T) {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	// Create and commit file foo1 in branch1
	repo.git("checkout", "-b", "branch1")
	repo.writeFile("foo1", "First line")
	repo.git("add", "foo1")
	repo.git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	repo.git("checkout", "-b", "branch2", "master")
	repo.writeFile("foo2", "First line")
	repo.git("add", "foo2")
	repo.git("commit", "-m\"\"")

	// Create and commit file foo3 in branch3
	repo.git("checkout", "-b", "branch3", "master")
	repo.writeFile("foo3", "First line")
	repo.git("add", "foo3")
	repo.git("commit", "-m\"\"")

	// Merge the 3 branches in master
	repo.git("checkout", "master")

	err := mainWithArgs(repo, "branch*")

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(repo.path, "foo1"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(repo.path, "foo2"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(repo.path, "foo3"))

	assert.Nil(t, err)
}
