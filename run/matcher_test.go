package run

import (
	"bytes"
	"github.com/lesfurets/git-octopus/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupRepo() (*OctopusContext, *bytes.Buffer) {
	context, out := CreateTestContext()
	head, _ := context.Repo.Git("rev-parse", "HEAD")
	context.Repo.Git("update-ref", "refs/heads/test1", head)
	context.Repo.Git("update-ref", "refs/remotes/origin/test1", head)
	context.Repo.Git("update-ref", "refs/remotes/origin/test2", head)
	context.Repo.Git("update-ref", "refs/remotes/origin/test3", head)

	return context, out
}

func TestResolveBranchListSimple(t *testing.T) {
	context, out := setupRepo()
	defer test.Cleanup(context.Repo)

	head, _ := context.Repo.Git("rev-parse", "HEAD")

	branchList := resolveBranchList(context.Repo, context.Logger, []string{"refs/heads/*"}, nil)

	expected := map[string]string{
		"refs/heads/master": head,
		"refs/heads/test1":  head,
	}

	assert.Equal(t, expected, branchList)
	outputString := out.String()
	assert.Contains(t, outputString,
		"2 branches (I)ncluded (2 matching, 0 (E)xcluded):\n"+
			"I  refs/heads/master\n"+
			"I  refs/heads/test1\n")
}

func TestResolveBranchListExclusion(t *testing.T) {
	context, out := setupRepo()
	defer test.Cleanup(context.Repo)

	head, _ := context.Repo.Git("rev-parse", "HEAD")

	branchList := resolveBranchList(context.Repo, context.Logger, []string{"refs/heads/*", "remotes/origin/*"}, []string{"*/test1"})

	expected := map[string]string{
		"refs/heads/master":         head,
		"refs/remotes/origin/test2": head,
		"refs/remotes/origin/test3": head,
	}

	assert.Equal(t, expected, branchList)

	outputString := out.String()
	assert.Contains(t, outputString,
		"3 branches (I)ncluded (5 matching, 2 (E)xcluded):\n"+
			"I  refs/heads/master\n"+
			"E  refs/heads/test1\n"+
			"E  refs/remotes/origin/test1\n"+
			"I  refs/remotes/origin/test2\n"+
			"I  refs/remotes/origin/test3\n")
}
