package main

import (
	"fmt"
	"os"
)

func main() {
	mainWithArgs(".", os.Args[1:])
}

func mainWithArgs(pwd string, args []string) {

	octopusConfig := getOctopusConfig(pwd, args)

	if octopusConfig.printVersion {
		fmt.Println("2.0")
		return
	}
}

func resolveBranchList(pwd string, patterns []string, excludedPatterns []string) map[string]string {
	result :=   parseLsRemote(git(pwd,  append([]string{"ls-remote", "."}, patterns...)...))

	if len(excludedPatterns) == 0 {
		return result
	}
	
	excludedRefs := parseLsRemote(git(pwd,  append([]string{"ls-remote", "."}, excludedPatterns...)...))
	for excludedRef, _ := range excludedRefs {
		delete(result, excludedRef)
	}

	return result
}