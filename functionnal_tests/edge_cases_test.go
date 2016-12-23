package functionnal_tests

import (
	"testing"
	"lesfurets/git-octopus/run"
	"github.com/stretchr/testify/assert"
	"lesfurets/git-octopus/test"
)

func TestVersionShort(t *testing.T) {
	context, out := run.CreateTestContext()
	defer test.Cleanup(context.Repo)

	run.Run(context, "-v")

	assert.Equal(t, "2.0\n", out.String())
}

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