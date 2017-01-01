package main

import (
	"github.com/lesfurets/git-octopus/git"
	"github.com/lesfurets/git-octopus/run"
	"log"
	"os"
	"os/signal"
)

func main() {
	repo := git.Repository{Path: "."}

	context := run.OctopusContext{
		Repo:   &repo,
		Logger: log.New(os.Stdout, "", 0),
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	go handleSignals(signalChan, &context)

	err := run.Run(&context, os.Args[1:]...)

	if err != nil {
		if len(err.Error()) > 0 {
			log.Fatalln(err.Error())
		}
		os.Exit(1)
	}
}

func handleSignals(signalChan chan os.Signal, context *run.OctopusContext) {
	initialHeadCommit, _ := context.Repo.Git("rev-parse", "HEAD")
	/*
	 The behavior of this is quite tricky. The signal is not only received on signalChan
	 but sent to subprocesses started by exec.Command as well. It is likely that
	 the main go routine is running one of those subprocess which will stop and return an error.
	 The error is handled by the Run function as any other error depending on where the execution was.

	 In the mean time, this routine is resetting the repo.

	 This is definitly an approximation that works in most cases.
	 */
	signal := <-signalChan
	context.Logger.Printf("Signal %v\n", signal.String())
	context.Repo.Git("reset", "-q", "--hard", initialHeadCommit)
	context.Repo.Git("clean", "-fd")
	os.Exit(1)
}