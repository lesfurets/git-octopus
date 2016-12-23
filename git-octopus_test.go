package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"lesfurets/git-octopus/git"
)

func createTestContext() (*octopusContext, *bytes.Buffer) {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")

	repo := git.Repository{Path: dir}

	repo.Git("init")
	repo.Git("commit", "--allow-empty", "-m\"first commit\"")

	out := bytes.NewBufferString("")

	context := octopusContext{
		repo:   &repo,
		logger: log.New(out, "", 0),
	}

	return &context, out
}

func cleanup(context *octopusContext) error {
	return os.RemoveAll(context.repo.Path)
}

func writeFile(repo *git.Repository, name string, lines ...string) {
	fileName := filepath.Join(repo.Path, name)
	ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

func TestVersionShort(t *testing.T) {
	context, out := createTestContext()
	defer cleanup(context)

	run(context, "-v")

	assert.Equal(t, "2.0\n", out.String())
}

func TestOctopusCommitConfigError(t *testing.T) {
	context, _ := createTestContext()
	defer cleanup(context)

	context.repo.Git("config", "octopus.commit", "bad_value")

	err := run(context, "-v")

	assert.NotNil(t, err)
}

func TestOctopusNoPatternGiven(t *testing.T) {
	context, out := createTestContext()
	defer cleanup(context)

	run(context)

	assert.Equal(t, "Nothing to merge. No pattern given\n", out.String())
}

func TestOctopusNoBranchMatching(t *testing.T) {
	context, out := createTestContext()
	defer cleanup(context)

	run(context, "refs/remotes/dumb/*", "refs/remotes/dumber/*")

	assert.Equal(t, "No branch matching \"refs/remotes/dumb/* refs/remotes/dumber/*\" were found\n", out.String())
}

func TestOctopusAlreadyUpToDate(t *testing.T) {
	context, _ := createTestContext()
	defer cleanup(context)

	writeFile(context.repo, "foo", "First line")
	context.repo.Git("add", "foo")
	context.repo.Git("commit", "-m\"first commit\"")
	// Create a branch on this first commit
	context.repo.Git("branch", "outdated_branch")
	expected, _ := context.repo.Git("rev-parse", "HEAD")

	run(context, "outdated_branch")

	actual, _ := context.repo.Git("rev-parse", "HEAD")

	assert.Equal(t, expected, actual)
}

func TestOctopus3branches(t *testing.T) {
	context, _ := createTestContext()
	defer cleanup(context)

	// Create and commit file foo1 in branch1
	context.repo.Git("checkout", "-b", "branch1")
	writeFile(context.repo, "foo1", "First line")
	context.repo.Git("add", "foo1")
	context.repo.Git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	context.repo.Git("checkout", "-b", "branch2", "master")
	writeFile(context.repo, "foo2", "First line")
	context.repo.Git("add", "foo2")
	context.repo.Git("commit", "-m\"\"")

	// Create and commit file foo3 in branch3
	context.repo.Git("checkout", "-b", "branch3", "master")
	writeFile(context.repo, "foo3", "First line")
	context.repo.Git("add", "foo3")
	context.repo.Git("commit", "-m\"\"")

	// Merge the 3 branches in master
	context.repo.Git("checkout", "master")

	err := run(context, "branch*")

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.repo.Path, "foo1"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.repo.Path, "foo2"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.repo.Path, "foo3"))

	assert.Nil(t, err)
}
