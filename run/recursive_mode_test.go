package run

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lesfurets/git-octopus/test"
	"github.com/stretchr/testify/assert"
)

// Basic merge of 3 branches. Asserts the resulting tree and the merge commit
func TestOctopus3BranchesRecursive(t *testing.T) {
	context, _ := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	// Create and commit file foo1 in branch1
	repo.Git("checkout", "-b", "branch1")
	writeFile(repo, "foo1", "First line")
	repo.Git("add", "foo1")
	repo.Git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	repo.Git("checkout", "-b", "branch2", "master")
	writeFile(repo, "foo2", "First line")
	repo.Git("add", "foo2")
	repo.Git("commit", "-m\"\"")

	// Create and commit file foo3 in branch3
	repo.Git("checkout", "-b", "branch3", "master")
	writeFile(repo, "foo3", "First line")
	repo.Git("add", "foo3")
	repo.Git("commit", "-m\"\"")

	// Merge the 3 branches in a new octopus branch
	repo.Git("checkout", "-b", "octopus", "master")

	err := Run(context, "-r", "branch*")
	assert.Nil(t, err)

	// The working tree should have the 3 files and status should be clean
	_, err = os.Open(filepath.Join(repo.Path, "foo1"))
	assert.Nil(t, err)
	_, err = os.Open(filepath.Join(repo.Path, "foo2"))
	assert.Nil(t, err)
	_, err = os.Open(filepath.Join(repo.Path, "foo3"))
	assert.Nil(t, err)

	status, _ := repo.Git("status", "--porcelain")
	assert.Empty(t, status)

	// octopus branch should contain the 3 branches
	_, err = repo.Git("merge-base", "--is-ancestor", "branch1", "octopus")
	assert.Nil(t, err)
	_, err = repo.Git("merge-base", "--is-ancestor", "branch2", "octopus")
	assert.Nil(t, err)
	_, err = repo.Git("merge-base", "--is-ancestor", "branch3", "octopus")
	assert.Nil(t, err)
}

func TestOctopus2BranchesRecursiveUnresolvedConflict(t *testing.T) {
	context, _ := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	// Create and commit file foo1 in branch1
	repo.Git("checkout", "-b", "branch1")
	writeFile(repo, "foo1", "First line from b1")
	repo.Git("add", "foo1")
	repo.Git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	repo.Git("checkout", "-b", "branch2", "master")
	writeFile(repo, "foo1", "First line from b2")
	repo.Git("add", "foo1")
	repo.Git("commit", "-m\"\"")

	// Merge the 3 branches in a new octopus branch
	repo.Git("checkout", "-b", "octopus", "master")

	err := Run(context, "-r", "branch*")
	assert.EqualError(t, err, "Unresolved merge conflict:\nAA foo1", "There should be a conflict")
}

func TestOctopus2BranchesRecursivePreRecordedConflict(t *testing.T) {
	context, _ := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	// Create and commit file foo1 in branch1
	repo.Git("checkout", "-b", "branch1")
	writeFile(repo, "foo1", "First line from b1")
	repo.Git("add", "foo1")
	repo.Git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	repo.Git("checkout", "-b", "branch2", "master")
	writeFile(repo, "foo1", "First line from b2")
	repo.Git("add", "foo1")
	repo.Git("commit", "-m\"\"")

	repo.Git("checkout", "-b", "rereretrain", "master")
	repo.Git("config", "--local", "rerere.enabled", "true")
	repo.Git("merge", "branch1")
	repo.Git("merge", "branch2")
	writeFile(repo, "foo1", "First line from b1\nFirst line from b2")
	repo.Git("add", "foo1")
	repo.Git("commit", "--no-edit")

	// Merge the 3 branches in a new octopus branch
	repo.Git("checkout", "-b", "octopus", "master")

	err := Run(context, "-r", "branch*")
	assert.Nil(t, err)
}

func TestOctopus2BranchesRecursiveFallbackPreRecordedConflict(t *testing.T) {
	context, _ := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	// Create and commit file foo1 in branch1
	repo.Git("checkout", "-b", "branch1")
	writeFile(repo, "foo1", "First line from b1")
	repo.Git("add", "foo1")
	repo.Git("commit", "-m\"\"")

	// Create and commit file foo2 in branch2
	repo.Git("checkout", "-b", "branch2", "master")
	writeFile(repo, "foo1", "First line from b2")
	repo.Git("add", "foo1")
	repo.Git("commit", "-m\"\"")

	repo.Git("checkout", "-b", "rereretrain", "master")
	repo.Git("config", "--local", "rerere.enabled", "true")
	repo.Git("merge", "branch1")
	repo.Git("merge", "branch2")
	writeFile(repo, "foo1", "First line from b1\nFirst line from b2")
	repo.Git("add", "foo1")
	repo.Git("commit", "--no-edit")

	// Merge the 3 branches in a new octopus branch
	repo.Git("checkout", "-b", "octopus", "master")

	err := Run(context, "-r", "-s=2", "branch*")
	assert.Nil(t, err)
}
