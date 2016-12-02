package main

import "fmt"
import "flag"

var printVersion bool

func main() {
	flag.BoolVar(&printVersion, "v", false, "prints the version of git-octopus")

	flag.Parse()

	mainWithArgs(printVersion)
}

func mainWithArgs(printVersion bool) {
	if printVersion {
		fmt.Printf("2.0")
	}
}
