package main

import (
	"github.com/lesfurets/git-octopus/git"
	"github.com/lesfurets/git-octopus/run"
	"log"
	"os"
)

func main() {
	repo := git.Repository{Path: "."}

	context := run.OctopusContext{
		Repo:   &repo,
		Logger: log.New(os.Stdout, "", 0),
	}

	err := run.Run(&context, os.Args[1:]...)

	if err != nil {
		if len(err.Error()) > 0 {
			log.Fatalln(err.Error())
		}
		os.Exit(1)
	}
}
