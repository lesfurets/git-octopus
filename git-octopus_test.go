package main

import "os"

func ExampleVersionShort() {
	os.Args = []string{"git-octopus", "-v"}
	mainWithArgs(true)
	// Output: 2.0
}
