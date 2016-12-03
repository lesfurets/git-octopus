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
		fmt.Printf("2.0\n")
		return
	}
}
