package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/lesfurets/git-octopus/git"
)

func main() {
	repo := git.Repository{Path: "."}
	mainWithArgs(&repo, os.Args[1:]...)
}

func mainWithArgs(repo *git.Repository, args ...string) {

	octopusConfig := getOctopusConfig(repo, args)

	if octopusConfig.printVersion {
		fmt.Println("2.0")
		return
	}

	if len(octopusConfig.patterns) == 0 {
		fmt.Println("Nothing to merge. No pattern given")
		return
	}

	branchList := resolveBranchList(repo, octopusConfig.patterns, octopusConfig.excludedPatterns)

	if len(branchList) == 0 {
		fmt.Printf("No branch matching \"%v\" were found\n", strings.Join(octopusConfig.patterns, " "))
		return
	}
}
