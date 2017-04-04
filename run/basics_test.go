package run

import (
	"github.com/lesfurets/git-octopus/git"
	"github.com/lesfurets/git-octopus/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	context, out := CreateTestContext()
	defer test.Cleanup(context.Repo)

	Run(context, "-v")

	assert.Equal(t, "2.0\n", out.String())
}

func writeFile(repo *git.Repository, name string, lines ...string) {
	fileName := filepath.Join(repo.Path, name)
	ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

func readFile(repo *git.Repository, name string) (string, error){
	fileName := filepath.Join(repo.Path, name)
	bytes, err := ioutil.ReadFile(fileName)
	return string(bytes), err
}

// Basic merge of 3 branches. Asserts the resulting tree and the merge commit
func TestOctopus3branches(t *testing.T) {
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

	err := Run(context, "branch*")
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

	// Assert the commit message
	commitMessage, _ := repo.Git("show", "--pretty=format:%B") // gets the commit body only

	assert.Contains(t, commitMessage,
		"Merged branches:\n"+
			"refs/heads/branch1\n"+
			"refs/heads/branch2\n"+
			"refs/heads/branch3\n"+
			"\nCommit created by git-octopus "+VERSION+".")
}

func TestOctopusDuplicateParent(t *testing.T) {
	context, outputBuffer := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	// Commit a new file on master
	writeFile(repo, "testFile", "line1")
	repo.Git("add", "testFile")
	repo.Git("commit", "-m", "Add first file")

	// Create an empty branch
	repo.Git("branch", "feature1")

	// Create a 2nd feature branch
	repo.Git("checkout", "-b", "feature2")
	// Edit the previously created file
	writeFile(repo, "testFile", "line1", "line2")
	repo.Git("add", "testFile")
	repo.Git("commit", "-m", "Add line2")

	// Run octopus from master
	repo.Git("checkout", "master")
	err := Run(context, "feature*")
	assert.Nil(t, err)

	// The file should be in the state of feature2 branch
	testFile, _ := readFile(repo, "testFile")
	assert.Equal(t, "line1\nline2", testFile)

	// The test case should reproduced the duplace parent error
	assert.Contains(t, outputBuffer.String(), "duplicate parent")
}