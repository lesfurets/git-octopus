package main

import (
	"reflect"
	"testing"
)

func TestResolveBranchListSimple(t *testing.T) {
	repo := initRepo()
	defer cleanupTestRepo(repo)

	head := repo.git("rev-parse", "HEAD")

	branchList := resolveBranchList(repo, []string{"refs/heads/*"}, nil)

	expected := map[string]string{
		"refs/heads/master": head,
		"refs/heads/test1":  head,
	}

	if !reflect.DeepEqual(branchList, expected) {
		t.Error(branchList)
	}
}

func TestResolveBranchListExclusion(t *testing.T) {
	repo := initRepo()
	defer cleanupTestRepo(repo)

	head := repo.git("rev-parse", "HEAD")

	branchList := resolveBranchList(repo, []string{"refs/heads/*", "remotes/origin/*"}, []string{"*/test1", "master"})

	expected := map[string]string{
		"refs/remotes/origin/test2": head,
	}

	if !reflect.DeepEqual(branchList, expected) {
		t.Error(branchList)
	}
}
