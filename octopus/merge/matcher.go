package merge

import (
	"bytes"
	"fmt"
	"github.com/lesfurets/git-octopus/octopus/git"
	"log"
	"strings"
)

func resolveBranchList(repo *git.Repository, logger *log.Logger, patterns []string, excludedPatterns []string) []git.LsRemoteEntry {
	lsRemote, _ := repo.Git(append([]string{"ls-remote", "."}, patterns...)...)
	includedRefs := git.ParseLsRemote(lsRemote)
	excludedRefs := []git.LsRemoteEntry{}

	totalCount := len(includedRefs)
	excludedCount := 0

	if len(excludedPatterns) > 0 {
		lsRemote, _ = repo.Git(append([]string{"ls-remote", "."}, excludedPatterns...)...)
		excludedRefs = git.ParseLsRemote(lsRemote)
	}

	tempBuffer := bytes.NewBufferString("")

	if totalCount == 0 {
		tempBuffer.WriteString(fmt.Sprintf("No branch matching \"%v\" were found\n", strings.Join(patterns, " ")))
	}

	result := []git.LsRemoteEntry{}

	for _, lsRemoteEntry := range includedRefs {
		excluded := false
		for _, excl := range excludedRefs {
			if excl.Ref == lsRemoteEntry.Ref {
				excludedCount++
				excluded = true
				break
			}
		}

		if excluded {
			tempBuffer.WriteString("E  ")
		} else {
			tempBuffer.WriteString("I  ")
			result = append(result, lsRemoteEntry)
		}
		tempBuffer.WriteString(lsRemoteEntry.Ref + "\n")
	}

	count := len(result)

	logger.Printf("%v branches (I)ncluded (%v matching, %v (E)xcluded):\n", count, totalCount, excludedCount)
	logger.Print(tempBuffer.String())

	return result
}
