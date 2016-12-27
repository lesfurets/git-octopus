package run

import (
	"bytes"
	"fmt"
	"github.com/lesfurets/git-octopus/git"
	"log"
	"strings"
)

func resolveBranchList(repo *git.Repository, logger *log.Logger, patterns []string, excludedPatterns []string) map[string]string {
	lsRemote, _ := repo.Git(append([]string{"ls-remote", "."}, patterns...)...)
	result := git.ParseLsRemote(lsRemote)
	excludedRefs := map[string]string{}

	totalCount := len(result)
	excludedCount := 0

	if len(excludedPatterns) > 0 {
		lsRemote, _ = repo.Git(append([]string{"ls-remote", "."}, excludedPatterns...)...)
		excludedRefs = git.ParseLsRemote(lsRemote)
	}

	tempBuffer := bytes.NewBufferString("")

	if totalCount == 0 {
		tempBuffer.WriteString(fmt.Sprintf("No branch matching \"%v\" were found\n", strings.Join(patterns, " ")))
	}

	for ref, _ := range result {
		_, ok := excludedRefs[ref]
		if ok {
			delete(result, ref)
			excludedCount++
			tempBuffer.WriteString("E  ")
		} else {
			tempBuffer.WriteString("I  ")
		}
		tempBuffer.WriteString(ref + "\n")
	}

	count := len(result)

	logger.Printf("%v branches (I)ncluded (%v matching, %v (E)xcluded):\n", count, totalCount, excludedCount)
	logger.Print(tempBuffer.String())

	return result
}
