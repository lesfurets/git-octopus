// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/lesfurets/git-octopus/git"
	"github.com/lesfurets/git-octopus/run"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: merge,
}

func merge(cmd *cobra.Command, args []string) {
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
	sig := <-signalChan
	context.Logger.Printf("Signal %v\n", sig.String())
	context.Repo.Git("reset", "-q", "--hard", initialHeadCommit)
	context.Repo.Git("clean", "-fd")
	os.Exit(1)
}

func init() {
	RootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
