package main

import (
	"reflect"
	"testing"
)

func ExampleVersionShort() {
	repo := createRepo()
	mainWithArgs(repo, []string{"-v"})
	// Output: 2.0
}

func TestResolveBranchListSimple(t *testing.T) {
	repo := createRepo()
	repo.run("git-octopus_test.sh")
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
	repo := createRepo()
	repo.run("git-octopus_test.sh")
	head := repo.git("rev-parse", "HEAD")

	branchList := resolveBranchList(repo, []string{"refs/heads/*", "remotes/origin/*"}, []string{"*/test1", "master"})

	expected := map[string]string{
		"refs/remotes/origin/test2": head,
	}

	if !reflect.DeepEqual(branchList, expected) {
		t.Error(branchList)
	}
}
