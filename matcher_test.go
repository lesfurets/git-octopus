package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupRepo() *repository {
	repo := createTestRepo()
	head, _ := repo.git("rev-parse", "HEAD")
	repo.git("update-ref", "refs/heads/test1", head)
	repo.git("update-ref", "refs/remotes/origin/test1", head)
	repo.git("update-ref", "refs/remotes/origin/test2", head)

	return repo
}

func TestResolveBranchListSimple(t *testing.T) {
	repo := setupRepo()
	defer cleanupTestRepo(repo)

	head, _ := repo.git("rev-parse", "HEAD")

	branchList := resolveBranchList(repo, []string{"refs/heads/*"}, nil)

	expected := map[string]string{
		"refs/heads/master": head,
		"refs/heads/test1":  head,
	}

	assert.Equal(t, expected, branchList)
}

func TestResolveBranchListExclusion(t *testing.T) {
	repo := setupRepo()
	defer cleanupTestRepo(repo)

	head, _ := repo.git("rev-parse", "HEAD")

	branchList := resolveBranchList(repo, []string{"refs/heads/*", "remotes/origin/*"}, []string{"*/test1", "master"})

	expected := map[string]string{
		"refs/remotes/origin/test2": head,
	}

	assert.Equal(t, expected, branchList)
}
