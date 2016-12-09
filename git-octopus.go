package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	repo := repository{path: "."}

	err := mainWithArgs(&repo, os.Args[1:]...)

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func mainWithArgs(repo *repository, args ...string) error {

	octopusConfig, err := getOctopusConfig(repo, args)

	if err != nil {
		return err
	}

	if octopusConfig.printVersion {
		fmt.Println("2.0")
		return nil
	}

	if len(octopusConfig.patterns) == 0 {
		fmt.Println("Nothing to merge. No pattern given")
		return nil
	}

	branchList := resolveBranchList(repo, octopusConfig.patterns, octopusConfig.excludedPatterns)

	if len(branchList) == 0 {
		fmt.Printf("No branch matching \"%v\" were found\n", strings.Join(octopusConfig.patterns, " "))
		return nil
	}

	return nil
}
