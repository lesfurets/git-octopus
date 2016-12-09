package main

func resolveBranchList(repo *repository, patterns []string, excludedPatterns []string) map[string]string {
	lsRemote, _ := repo.git(append([]string{"ls-remote", "."}, patterns...)...)
	result := parseLsRemote(lsRemote)

	if len(excludedPatterns) == 0 {
		return result
	}

	lsRemote, _ = repo.git(append([]string{"ls-remote", "."}, excludedPatterns...)...)
	excludedRefs := parseLsRemote(lsRemote)
	for excludedRef, _ := range excludedRefs {
		delete(result, excludedRef)
	}

	return result
}
