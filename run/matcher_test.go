package run

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupRepo() *OctopusContext {
	context, _ := createTestContext()
	head, _ := context.Repo.Git("rev-parse", "HEAD")
	context.Repo.Git("update-ref", "refs/heads/test1", head)
	context.Repo.Git("update-ref", "refs/remotes/origin/test1", head)
	context.Repo.Git("update-ref", "refs/remotes/origin/test2", head)

	return context
}

func TestResolveBranchListSimple(t *testing.T) {
	context := setupRepo()
	defer cleanup(context)

	head, _ := context.Repo.Git("rev-parse", "HEAD")

	branchList := resolveBranchList(context.Repo, []string{"refs/heads/*"}, nil)

	expected := map[string]string{
		"refs/heads/master": head,
		"refs/heads/test1":  head,
	}

	assert.Equal(t, expected, branchList)
}

func TestResolveBranchListExclusion(t *testing.T) {
	context := setupRepo()
	defer cleanup(context)

	head, _ := context.Repo.Git("rev-parse", "HEAD")

	branchList := resolveBranchList(context.Repo, []string{"refs/heads/*", "remotes/origin/*"}, []string{"*/test1", "master"})

	expected := map[string]string{
		"refs/remotes/origin/test2": head,
	}

	assert.Equal(t, expected, branchList)
}
