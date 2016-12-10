package main

func resolveBranchList(repo *repository, patterns []string, excludedPatterns []string) map[string]string {
	result := parseLsRemote(repo.git(append([]string{"ls-remote", "."}, patterns...)...))

	if len(excludedPatterns) == 0 {
		return result
	}

	excludedRefs := parseLsRemote(repo.git(append([]string{"ls-remote", "."}, excludedPatterns...)...))
	for excludedRef, _ := range excludedRefs {
		delete(result, excludedRef)
	}

	return result
}
