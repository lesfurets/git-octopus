package run

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"strings"

	"github.com/lesfurets/git-octopus/config"
	"github.com/lesfurets/git-octopus/git"
)

type OctopusContext struct {
	Repo   *git.Repository
	Logger *log.Logger
}

const VERSION = "2.0"

func Run(context *OctopusContext, args ...string) error {

	octopusConfig, err := config.GetOctopusConfig(context.Repo, args)

	if err != nil {
		return err
	}

	if octopusConfig.PrintVersion {
		context.Logger.Println(VERSION)
		return nil
	}

	if len(octopusConfig.Patterns) == 0 {
		context.Logger.Println("Nothing to merge. No pattern given")
		return nil
	}

	status, _ := context.Repo.Git("status", "--porcelain")

	// This is not formally required but it would be an ambiguous behaviour to let git-octopus run on unclean state.
	if len(status) != 0 {
		return errors.New("The repository has to be clean.")
	}

	branchList := resolveBranchList(context.Repo, context.Logger, octopusConfig.Patterns, octopusConfig.ExcludedPatterns)

	if len(branchList) == 0 {
		return nil
	}

	context.Logger.Println()

	mergeStrategy := chooseMergeStrategy(octopusConfig)
	err = mergeStrategy(context, octopusConfig, branchList)

	context.Logger.Println()

	return err
}

type strategy func(context *OctopusContext, octopusConfig *config.OctopusConfig, branchList []git.LsRemoteEntry) error

func chooseMergeStrategy(octopusConfig *config.OctopusConfig) strategy {
	chunkMode := octopusConfig.ChunkSize > 0
	if octopusConfig.RecursiveMode {
		if chunkMode {
			return chunckBranches(octopusWithRecursiveFallbackStrategy)
		}
		return recursiveStrategy
	}
	if chunkMode {
		return chunckBranches(octopusStrategy)
	}
	return octopusStrategy
}

func chunckBranches(mergeStrategy strategy) strategy {
	return func(context *OctopusContext, octopusConfig *config.OctopusConfig, branchList []git.LsRemoteEntry) error {
		var remaning []git.LsRemoteEntry = branchList
		chunkSize := octopusConfig.ChunkSize
		acc := 1
		lenBranchList := len(branchList)
		context.Logger.Printf("Will merge %d branches by chunks of %d", lenBranchList, chunkSize)
		for len(remaning) > 0 {
			var current []git.LsRemoteEntry
			if len(remaning) > chunkSize {
				current, remaning = remaning[:chunkSize], remaning[chunkSize:]
			} else {
				current, remaning = remaning, nil
			}
			lcur := len(current)
			context.Logger.Printf("Merging chunks %d to %d (out of %d)\n", acc, acc+lcur-1, lenBranchList)
			acc += lcur
			err := mergeStrategy(context, octopusConfig, current)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func octopusWithRecursiveFallbackStrategy(context *OctopusContext, octopusConfig *config.OctopusConfig,
	branchList []git.LsRemoteEntry) error {

	currentHeadCommit, _ := context.Repo.Git("rev-parse", "HEAD")
	err := octopusStrategy(context, octopusConfig, branchList)
	if err != nil {
		if len(err.Error()) > 0 {
			context.Logger.Println(err.Error())
		}
		context.Logger.Printf("Octopus strategy failed for branches %v, fallback to one by one recursive merge\n", branchList)
		context.Repo.Git("reset", "-q", "--hard", currentHeadCommit)
		err = recursiveStrategy(context, octopusConfig, branchList)
	}
	return err
}

func octopusStrategy(context *OctopusContext, octopusConfig *config.OctopusConfig,
	branchList []git.LsRemoteEntry) error {
	initialHeadCommit, _ := context.Repo.Git("rev-parse", "HEAD")
	parents, err := mergeHeads(context, branchList)

	if !octopusConfig.DoCommit {
		context.Repo.Git("reset", "-q", "--hard", initialHeadCommit)
	}

	if err != nil {
		return err
	}

	// parents always contains HEAD. We need at lease 2 parents to create a merge commit
	if octopusConfig.DoCommit && parents != nil && len(parents) > 1 {
		tree, _ := context.Repo.Git("write-tree")
		args := []string{"commit-tree"}
		for _, parent := range parents {
			args = append(args, "-p", parent)
		}
		args = append(args, "-m", octopusCommitMessage(branchList), tree)
		commit, _ := context.Repo.Git(args...)
		context.Repo.Git("update-ref", "HEAD", commit)
	}
	return nil
}

func recursiveStrategy(context *OctopusContext, octopusConfig *config.OctopusConfig,
	branchList []git.LsRemoteEntry) error {
	context.Logger.Println("Merging using recursive mode")
	initialHeadCommit, _ := context.Repo.Git("rev-parse", "HEAD")
	_, err := mergeRecursive(context, branchList)

	if !octopusConfig.DoCommit {
		context.Repo.Git("reset", "-q", "--hard", initialHeadCommit)
	}

	return err
}

func mergeRecursive(context *OctopusContext, remotes []git.LsRemoteEntry) ([]string, error) {
	head, _ := context.Repo.Git("rev-parse", "--verify", "-q", "HEAD")
	mrc := []string{head}
	for _, lsRemoteEntry := range remotes {
		context.Logger.Println("Merging " + lsRemoteEntry.Ref)
		log, _ := context.Repo.Git("merge", "--no-commit", "--rerere-autoupdate", lsRemoteEntry.Ref)
		if len(log) > 0 {
			context.Logger.Println(log)
		}

		status, _ := context.Repo.Git("status", "--porcelain")
		if isMergeStatusOk(context, status) {
			context.Repo.Git("commit", "--no-edit")
			mrc = append(mrc, lsRemoteEntry.Sha1)
		} else {
			return nil, errors.New("Unresolved merge conflict:\n" + status)
		}
	}
	context.Repo.Git("commit", "-m", octopusCommitMessage(remotes), "--allow-empty")
	return mrc, nil
}

// Takes the output of git-ls-remote. Returns a map refsname => sha1
func isMergeStatusOk(context *OctopusContext, status string) bool {
	scanner := bufio.NewScanner(strings.NewReader(status))
	for scanner.Scan() {
		split := strings.Fields(scanner.Text())

		switch split[0] {
		case "DD", "AU", "UD", "UA", "DU", "AA", "UU":
			return false
		}
	}

	return true
}

// The logic of this function is copied directly from git-merge-octopus.sh
func mergeHeads(context *OctopusContext, remotes []git.LsRemoteEntry) ([]string, error) {
	head, _ := context.Repo.Git("rev-parse", "--verify", "-q", "HEAD")

	mrc := []string{head}
	mrt, _ := context.Repo.Git("write-tree")
	nonFfMerge := false

	for _, lsRemoteEntry := range remotes {

		common, err := context.Repo.Git(append([]string{"merge-base", "--all", lsRemoteEntry.Sha1}, mrc...)...)

		if err != nil {
			return nil, errors.New("Unable to find common commit with " + lsRemoteEntry.Ref)
		}

		if common == lsRemoteEntry.Sha1 {
			context.Logger.Println("Already up-to-date with " + lsRemoteEntry.Ref)
			continue
		}

		if len(mrc) == 1 && common == mrc[0] && !nonFfMerge {
			context.Logger.Println("Fast-forwarding to: " + lsRemoteEntry.Ref)
			_, err := context.Repo.Git("read-tree", "-u", "-m", head, lsRemoteEntry.Sha1)

			if err != nil {
				return nil, nil
			}

			mrc[0] = lsRemoteEntry.Sha1
			mrt, _ = context.Repo.Git("write-tree")
			continue
		}

		nonFfMerge = true

		context.Logger.Println("Trying simple merge with " + lsRemoteEntry.Ref)

		commonArray := strings.Split(common, "\n")
		_, err = context.Repo.Git(append([]string{"read-tree", "-u", "-m", "--aggressive"},
			append(commonArray, mrt, lsRemoteEntry.Sha1)...)...)

		if err != nil {
			return nil, err
		}

		next, err := context.Repo.Git("write-tree")

		if err != nil {
			context.Logger.Println("Simple merge did not work, trying automatic merge.")
			_, err = context.Repo.Git("merge-index", "-o", "git-merge-one-file", "-a")

			if err != nil {
				context.Logger.Println("Automated merge did not work.")
				context.Logger.Println("Should not be doing an Octopus.")
				return nil, errors.New("")
			}

			next, _ = context.Repo.Git("write-tree")
		}

		mrc = append(mrc, lsRemoteEntry.Sha1)
		mrt = next
	}

	return mrc, nil
}

func octopusCommitMessage(remotes []git.LsRemoteEntry) string {
	buf := bytes.NewBufferString("Merged branches:\n")
	for _, lsRemoteEntry := range remotes {
		buf.WriteString(lsRemoteEntry.Ref + "\n")
	}
	buf.WriteString("\nCommit created by git-octopus " + VERSION + ".\n")
	return buf.String()
}
