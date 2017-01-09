package run

import (
	"github.com/lesfurets/git-octopus/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFastForward(t *testing.T) {
	context, _ := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	// The repo is on master branch with an empty tree
	// Create a branch with a new file
	repo.Git("checkout", "-b", "new_branch")
	writeFile(repo, "foo", "bar")
	repo.Git("add", "foo")
	repo.Git("commit", "-m", "added foo")

	repo.Git("checkout", "master")

	expected, _ := repo.Git("rev-parse", "HEAD")

	Run(context, "-n", "new_branch")

	actual, _ := repo.Git("rev-parse", "HEAD")
	assert.Equal(t, expected, actual)

	status, _ := repo.Git("status", "--porcelain")
	assert.Empty(t, status)
}

func TestConflictState(t *testing.T) {
	context, _ := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	writeFile(repo, "foo", "line 1", "")
	repo.Git("add", ".")
	repo.Git("commit", "-m", "added foo")

	writeFile(repo, "foo", "line 1", "line 2")
	repo.Git("commit", "-a", "-m", "edited foo")

	repo.Git("checkout", "-b", "a_branch", "HEAD^")

	writeFile(repo, "foo", "line 1", "line 2 bis")
	repo.Git("commit", "-a", "-m", "edited foo in parallel to master")

	repo.Git("checkout", "master")
	expected, _ := repo.Git("rev-parse", "HEAD")

	err := Run(context, "-n", "a_branch")

	assert.NotNil(t, err)
	actual, _ := repo.Git("rev-parse", "HEAD")
	assert.Equal(t, expected, actual)

	status, _ := repo.Git("status", "--porcelain")
	assert.Empty(t, status)
}
