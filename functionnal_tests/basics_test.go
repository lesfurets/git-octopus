package functionnal_tests

import (
	"testing"
	"lesfurets/git-octopus/run"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
	"lesfurets/git-octopus/git"
	"lesfurets/git-octopus/test"
)

func writeFile(repo *git.Repository, name string, lines ...string) {
	fileName := filepath.Join(repo.Path, name)
	ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

func TestOctopusAlreadyUpToDate(t *testing.T) {
	context, _ := run.CreateTestContext()
	defer test.Cleanup(context.Repo)

	writeFile(context.Repo, "foo", "First line")
	context.Repo.Git("add", "foo")
	context.Repo.Git("commit", "-m\"first commit\"")
	// Create a branch on this first commit
	context.Repo.Git("branch", "outdated_branch")
	expected, _ := context.Repo.Git("rev-parse", "HEAD")

	run.Run(context, "outdated_branch")

	actual, _ := context.Repo.Git("rev-parse", "HEAD")

	assert.Equal(t, expected, actual)
}

func TestOctopus3branches(t *testing.T) {
	context, _ := run.CreateTestContext()
	defer test.Cleanup(context.Repo)

	// Create and commit file foo1 in branch1
	context.Repo.Git("checkout", "-b", "branch1")
	writeFile(context.Repo, "foo1", "First line")
	context.Repo.Git("add", "foo1")
	context.Repo.Git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	context.Repo.Git("checkout", "-b", "branch2", "master")
	writeFile(context.Repo, "foo2", "First line")
	context.Repo.Git("add", "foo2")
	context.Repo.Git("commit", "-m\"\"")

	// Create and commit file foo3 in branch3
	context.Repo.Git("checkout", "-b", "branch3", "master")
	writeFile(context.Repo, "foo3", "First line")
	context.Repo.Git("add", "foo3")
	context.Repo.Git("commit", "-m\"\"")

	// Merge the 3 branches in master
	context.Repo.Git("checkout", "master")

	err := run.Run(context, "branch*")

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.Repo.Path, "foo1"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.Repo.Path, "foo2"))

	assert.Nil(t, err)

	_, err = os.Open(filepath.Join(context.Repo.Path, "foo3"))

	assert.Nil(t, err)
}
