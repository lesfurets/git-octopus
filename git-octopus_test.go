package main

import (
	"reflect"
	"testing"
)

func ExampleVersionShort() {
	repo := createRepo()
	mainWithArgs(repo, "-v")
	// Output: 2.0
}

func initRepo() *repository {
	repo := createRepo()
	repo.git("commit", "--allow-empty", "-m\"first commit\"")
	head := repo.git("rev-parse", "HEAD")
	repo.git("update-ref", "refs/heads/test1", head)
	repo.git("update-ref", "refs/remotes/origin/test1", head)
	repo.git("update-ref", "refs/remotes/origin/test2", head)

	return repo
}

func TestResolveBranchListSimple(t *testing.T) {
	repo := initRepo()
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
	head := repo.git("rev-parse", "HEAD")

	branchList := resolveBranchList(repo, []string{"refs/heads/*", "remotes/origin/*"}, []string{"*/test1", "master"})

	expected := map[string]string{
		"refs/remotes/origin/test2": head,
	}

	if !reflect.DeepEqual(branchList, expected) {
		t.Error(branchList)
	}
}

func ExampleOctopusNoPatternGiven() {
	repo := createRepo()
	mainWithArgs(repo)
	// Output: Nothing to merge. No pattern given
}

func ExampleOctopusNoBranchMatching() {
	repo := createRepo()
	mainWithArgs(repo, "refs/remotes/dumb/*", "refs/remotes/dumber/*")
	// Output: No branch matching "refs/remotes/dumb/* refs/remotes/dumber/*" were found
}