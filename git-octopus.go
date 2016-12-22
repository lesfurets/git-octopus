package main

import (
	"errors"
	"log"
	"os"
	"strings"
)

type octopusContext struct {
	repo   *repository
	logger *log.Logger
}

func main() {
	repo := repository{path: "."}

	context := octopusContext{
		repo:   &repo,
		logger: log.New(os.Stdout, "", 0),
	}

	err := run(&context, os.Args[1:]...)

	if err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
}

func run(context *octopusContext, args ...string) error {

	octopusConfig, err := getOctopusConfig(context.repo, args)

	if err != nil {
		return err
	}

	if octopusConfig.printVersion {
		context.logger.Println("2.0")
		return nil
	}

	if len(octopusConfig.patterns) == 0 {
		context.logger.Println("Nothing to merge. No pattern given")
		return nil
	}

	branchList := resolveBranchList(context.repo, octopusConfig.patterns, octopusConfig.excludedPatterns)

	if len(branchList) == 0 {
		context.logger.Printf("No branch matching \"%v\" were found\n", strings.Join(octopusConfig.patterns, " "))
		return nil
	}

	parents, err := mergeHeads(context, branchList)

	if err != nil {
		return err
	}

	if octopusConfig.doCommit {
		tree, _ := context.repo.git("write-tree")
		commit, _ := context.repo.git("commit-tree", "-p", strings.Join(parents, " -p "), "-m", octopusCommitMessage(branchList), tree)
		context.repo.git("update-ref", "HEAD", commit)
	}

	return nil
}

// The logic of this function is copied directly from git-merge-octopus.sh
func mergeHeads(context *octopusContext, remotes map[string]string) ([]string, error) {
	head, _ := context.repo.git("rev-parse", "--verify", "-q", "HEAD")

	alreadyUpToDate := true
	for _, sha1 := range remotes {
		_, err := context.repo.git("merge-base", "--is-ancestor", sha1, "HEAD")
		if err != nil {
			alreadyUpToDate = false
		}
	}
	// This prevents git-octopus to create a commit when there's nothing to merge,
	// i.e. no feature branches but only master.
	if alreadyUpToDate {
		context.logger.Println("Already up to date")
		return nil, nil
	}

	mrc := []string{head}
	mrt, _ := context.repo.git("write-tree")
	nonFfMerge := false

	for prettyRemoteName, sha1 := range remotes {

		common, err := context.repo.git(append([]string{"merge-base", "--all", sha1}, mrc...)...)

		if err != nil {
			return nil, errors.New("Unable to find common commit with " + prettyRemoteName)
		}

		if common == sha1 {
			context.logger.Println("Already up-to-date with " + prettyRemoteName)
			continue
		}

		if len(mrc) == 1 && common == mrc[0] && !nonFfMerge {
			context.logger.Println("Fast-forwarding to: " + prettyRemoteName)
			_, err := context.repo.git("read-tree", "-u", "-m", head, sha1)

			if err != nil {
				return nil, nil
			}

			mrc[0] = sha1
			mrt, _ = context.repo.git("write-tree")
			continue
		}

		nonFfMerge = true

		context.logger.Println("Trying simple merge with " + prettyRemoteName)

		_, err = context.repo.git("read-tree", "-u", "-m", "--aggressive", common, mrt, sha1)

		if err != nil {
			return nil, err
		}

		next, err := context.repo.git("write-tree")

		if err != nil {
			context.logger.Println("Simple merge did not work, trying automatic merge.")
			_, err = context.repo.git("merge-index", "-o", "git-merge-one-file", "-a")

			if err != nil {
				context.logger.Println("Automated merge did not work.")
				context.logger.Println("Should not be doing an Octopus.")
				return nil, errors.New("")
			}

			next, _ = context.repo.git("write-tree")
		}

		mrc = append(mrc, sha1)
		mrt = next
	}

	return mrc, nil
}

func octopusCommitMessage(remotes map[string]string) string {
	return "octopus commit"
}
