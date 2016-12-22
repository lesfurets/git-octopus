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
)

func createTestContext() (*octopusContext, *bytes.Buffer) {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")

	repo := repository{path: dir}

	repo.git("init")
	repo.git("commit", "--allow-empty", "-m\"first commit\"")

	out := bytes.NewBufferString("")

	context := octopusContext{
		repo:   &repo,
		logger: log.New(out, "", 0),
	}

	return &context, out
}

func cleanup(context *octopusContext) error {
	return os.RemoveAll(context.repo.path)
}

func (repo *repository) writeFile(name string, lines ...string) {
	fileName := filepath.Join(repo.path, name)
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

	context.repo.git("config", "octopus.commit", "bad_value")

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

	context.repo.writeFile("foo", "First line")
	context.repo.git("add", "foo")
	context.repo.git("commit", "-m\"first commit\"")
	// Create a branch on this first commit
	context.repo.git("branch", "outdated_branch")
	expected, _ := context.repo.git("rev-parse", "HEAD")

	run(context, "outdated_branch")

	actual, _ := context.repo.git("rev-parse", "HEAD")

	assert.Equal(t, expected, actual)
}

func TestOctopus3branches(t *testing.T) {
	context, _ := createTestContext()
	defer cleanup(context)

	// Create and commit file foo1 in branch1
	context.repo.git("checkout", "-b", "branch1")
	context.repo.writeFile("foo1", "First line")
	context.repo.git("add", "foo1")
	context.repo.git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	context.repo.git("checkout", "-b", "branch2", "master")
	context.repo.writeFile("foo2", "First line")
	context.repo.git("add", "foo2")
	context.repo.git("commit", "-m\"\"")

	// Create and commit file foo3 in branch3
	context.repo.git("checkout", "-b", "branch3", "master")
	context.repo.writeFile("foo3", "First line")
	context.repo.git("add", "foo3")
	context.repo.git("commit", "-m\"\"")

	// Merge the 3 branches in master
	context.repo.git("checkout", "master")

	err := run(context, "branch*")

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.repo.path, "foo1"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.repo.path, "foo2"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.repo.path, "foo3"))

	assert.Nil(t, err)
}
