package main

import "github.com/lesfurets/git-octopus/git"

func resolveBranchList(repo *git.Repository, patterns []string, excludedPatterns []string) map[string]string {
	result := git.ParseLsRemote(repo.Git(append([]string{"ls-remote", "."}, patterns...)...))

	if len(excludedPatterns) == 0 {
		return result
	}

	excludedRefs := git.ParseLsRemote(repo.Git(append([]string{"ls-remote", "."}, excludedPatterns...)...))
	for excludedRef, _ := range excludedRefs {
		delete(result, excludedRef)
	}

	return result
}
