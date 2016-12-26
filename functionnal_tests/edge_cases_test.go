package functionnal_tests

import (
	"testing"
	"lesfurets/git-octopus/run"
	"github.com/stretchr/testify/assert"
	"lesfurets/git-octopus/test"
)

func TestOctopusCommitConfigError(t *testing.T) {
	context, _ := run.CreateTestContext()
	defer test.Cleanup(context.Repo)

	context.Repo.Git("config", "octopus.commit", "bad_value")

	err := run.Run(context, "-v")

	assert.NotNil(t, err)
}

func TestOctopusNoPatternGiven(t *testing.T) {
	context, out := run.CreateTestContext()
	defer test.Cleanup(context.Repo)

	run.Run(context)

	assert.Equal(t, "Nothing to merge. No pattern given\n", out.String())
}

func TestOctopusNoBranchMatching(t *testing.T) {
	context, out := run.CreateTestContext()
	defer test.Cleanup(context.Repo)

	run.Run(context, "refs/remotes/dumb/*", "refs/remotes/dumber/*")

	assert.Equal(t, "No branch matching \"refs/remotes/dumb/* refs/remotes/dumber/*\" were found\n", out.String())
}

// Merge a branch that is already merged.
// Should be noop and print something accordingly
func TestOctopusAlreadyUpToDate(t *testing.T) {
	context, out := run.CreateTestContext()
	defer test.Cleanup(context.Repo)

	// commit a file in master
	writeFile(context.Repo, "foo", "First line")
	context.Repo.Git("add", "foo")
	context.Repo.Git("commit", "-m\"first commit\"")

	// Create a branch on this first commit.
	// master and outdated_branch are on the same commit
	context.Repo.Git("branch", "outdated_branch")

	expected, _ := context.Repo.Git("rev-parse", "HEAD")

	err := run.Run(context, "outdated_branch")

	actual, _ := context.Repo.Git("rev-parse", "HEAD")

	// HEAD should point to the same commit
	assert.Equal(t, expected, actual)

	// This is a normal behavious, no error should be raised
	assert.Nil(t, err)

	assert.Contains(t, out.String(), "Already up-to-date with refs/heads/outdated_branch")
}