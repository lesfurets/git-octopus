package merge

import (
	"bytes"
	"fmt"
	"github.com/lesfurets/git-octopus/octopus/git"
	"github.com/lesfurets/git-octopus/octopus/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func CreateTestContext() (*MergeContext, *bytes.Buffer) {
	dir := test.CreateTempDir()

	repo := git.Repository{Path: dir}

	repo.Git("init")
	repo.Git("config", "user.name", "gotest")
	repo.Git("config", "user.email", "gotest@golang.com")
	_, err := repo.Git("commit", "--allow-empty", "-m\"first commit\"")

	if err != nil {
		fmt.Println("There's something wrong with the git installation:")
		fmt.Println(err.Error())
	}

	out := bytes.NewBufferString("")

	context := MergeContext{
		Repo:   &repo,
		Logger: log.New(out, "", 0),
	}

	return &context, out
}

func writeFile(repo *git.Repository, name string, lines ...string) {
	fileName := filepath.Join(repo.Path, name)
	ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
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

	config := Config{
		Patterns: []string{"branch*"},
	}
	err := Merge(context, &config)
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
			"\nCommit created by git-octopus.")
}

func TestOctopusNoPatternGiven(t *testing.T) {
	context, out := CreateTestContext()
	defer test.Cleanup(context.Repo)

	config := Config{}
	Merge(context, &config)

	assert.Equal(t, "Nothing to merge. No pattern given\n", out.String())
}

func TestOctopusNoBranchMatching(t *testing.T) {
	context, out := CreateTestContext()
	defer test.Cleanup(context.Repo)

	config := Config{
		Patterns: []string{"refs/remotes/dumb/*", "refs/remotes/dumber/*"},
	}
	Merge(context, &config)

	assert.Contains(t, out.String(), "No branch matching \"refs/remotes/dumb/* refs/remotes/dumber/*\" were found\n")
}

// Merge a branch that is already merged.
// Should be noop and print something accordingly
func TestOctopusAlreadyUpToDate(t *testing.T) {
	context, out := CreateTestContext()
	defer test.Cleanup(context.Repo)

	// commit a file in master
	writeFile(context.Repo, "foo", "First line")
	context.Repo.Git("add", "foo")
	context.Repo.Git("commit", "-m\"first commit\"")

	// Create a branch on this first commit.
	// master and outdated_branch are on the same commit
	context.Repo.Git("branch", "outdated_branch")

	expected, _ := context.Repo.Git("rev-parse", "HEAD")

	config := Config{
		Patterns: []string{"outdated_branch"},
	}
	err := Merge(context, &config)

	actual, _ := context.Repo.Git("rev-parse", "HEAD")

	// HEAD should point to the same commit
	assert.Equal(t, expected, actual)

	// This is a normal behavious, no error should be raised
	assert.Nil(t, err)

	assert.Contains(t, out.String(), "Already up-to-date with refs/heads/outdated_branch")
}

// git-octopus should prevent from running if status is not clean
func TestUncleanStateFail(t *testing.T) {
	context, _ := CreateTestContext()
	defer test.Cleanup(context.Repo)

	// create and commit a file
	writeFile(context.Repo, "foo", "First line")

	config := Config{
		Patterns: []string{"*"},
	}
	err := Merge(context, &config)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "The repository has to be clean.")
	}
}

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

	config := Config{
		NoCommit: true,
		Patterns: []string{"new_branch"},
	}
	Merge(context, &config)

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

	config := Config{
		NoCommit: true,
		Patterns: []string{"a_branch"},
	}
	err := Merge(context, &config)

	assert.NotNil(t, err)
	actual, _ := repo.Git("rev-parse", "HEAD")
	assert.Equal(t, expected, actual)

	status, _ := repo.Git("status", "--porcelain")
	assert.Empty(t, status)
}
