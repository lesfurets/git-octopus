package run

import "lesfurets/git-octopus/git"

func resolveBranchList(repo *git.Repository, patterns []string, excludedPatterns []string) map[string]string {
	lsRemote, _ := repo.Git(append([]string{"ls-remote", "."}, patterns...)...)
	result := git.ParseLsRemote(lsRemote)

	if len(excludedPatterns) == 0 {
		return result
	}

	lsRemote, _ = repo.Git(append([]string{"ls-remote", "."}, excludedPatterns...)...)
	excludedRefs := git.ParseLsRemote(lsRemote)
	for excludedRef := range excludedRefs {
		delete(result, excludedRef)
	}

	return result
}
