package main

import (
	"log"
	"os"
	"lesfurets/git-octopus/git"
	"lesfurets/git-octopus/run"
)

func main() {
	repo := git.Repository{Path: "."}

	context := run.OctopusContext{
		Repo:   &repo,
		Logger: log.New(os.Stdout, "", 0),
	}

	err := run.Run(&context, os.Args[1:]...)

	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
}

