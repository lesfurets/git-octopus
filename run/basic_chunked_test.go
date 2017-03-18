package run

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lesfurets/git-octopus/test"
	"github.com/stretchr/testify/assert"
)

// Basic merge of 3 branches with chunk of 2. Asserts the resulting tree and the merge commit
func TestOctopus3BranchesChunked(t *testing.T) {
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

	err := Run(context, "-s=2", "branch*")
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

	// // Assert the commit message
	commitMessage1, _ := repo.Git("show", "--pretty=format:%B") // gets the commit body only
	assert.Contains(t, commitMessage1,
		"Merged branches:\n"+
			"refs/heads/branch3\n"+
			"\nCommit created by git-octopus "+VERSION+".")

	commitMessage2, _ := repo.Git("show", "--pretty=format:%B", "HEAD^") // gets the commit body only
	assert.Contains(t, commitMessage2,
		"Merged branches:\n"+
			"refs/heads/branch1\n"+
			"refs/heads/branch2\n"+
			"\nCommit created by git-octopus "+VERSION+".")
}
