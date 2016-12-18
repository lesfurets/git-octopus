package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	repo := repository{path: "."}

	err := mainWithArgs(&repo, os.Args[1:]...)

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func mainWithArgs(repo *repository, args ...string) error {

	octopusConfig, err := getOctopusConfig(repo, args)

	if err != nil {
		return err
	}

	if octopusConfig.printVersion {
		fmt.Println("2.0")
		return nil
	}

	if len(octopusConfig.patterns) == 0 {
		fmt.Println("Nothing to merge. No pattern given")
		return nil
	}

	branchList := resolveBranchList(repo, octopusConfig.patterns, octopusConfig.excludedPatterns)

	if len(branchList) == 0 {
		fmt.Printf("No branch matching \"%v\" were found\n", strings.Join(octopusConfig.patterns, " "))
		return nil
	}

	parents, err := mergeHeads(repo, branchList)

	if err != nil {
		return err
	}

	if octopusConfig.doCommit {
		tree, _ := repo.git("write-tree")
		commit, _ := repo.git("commit-tree", "-p", strings.Join(parents, " -p "), "-m", octopusCommitMessage(branchList), tree)
		repo.git("update-ref", "HEAD", commit)
	}

	return nil
}

// The logic of this function is copied directly from git-merge-octopus.sh
func mergeHeads(repo *repository, remotes map[string]string) ([]string, error) {
	head, _ := repo.git("rev-parse", "--verify", "-q", "HEAD")

	alreadyUpToDate := true
	for _, sha1 := range remotes {
		_, err := repo.git("merge-base", "--is-ancestor", sha1, "HEAD")
		if err != nil {
			alreadyUpToDate = false
		}
	}
	// This prevents git-octopus to create a commit when there's nothing to merge,
	// i.e. no feature branches but only master.
	if alreadyUpToDate {
		fmt.Println("Already up to date")
		return nil, nil
	}

	mrc := []string{head}
	mrt, _ := repo.git("write-tree")
	nonFfMerge := false

	for prettyRemoteName, sha1 := range remotes {

		common, err := repo.git(append([]string{"merge-base", "--all", sha1}, mrc...)...)

		if err != nil {
			return nil, errors.New("Unable to find common commit with " + prettyRemoteName)
		}

		if common == sha1 {
			fmt.Println("Already up-to-date with " + prettyRemoteName)
			continue
		}

		if len(mrc) == 1 && common == mrc[0] && !nonFfMerge {
			fmt.Println("Fast-forwarding to: " + prettyRemoteName)
			_, err := repo.git("read-tree", "-u", "-m", head, sha1)

			if err != nil {
				return nil, nil
			}

			mrc[0] = sha1
			mrt, _ = repo.git("write-tree")
			continue
		}

		nonFfMerge = true

		fmt.Println("Trying simple merge with " + prettyRemoteName)

		_, err = repo.git("read-tree", "-u", "-m", "--aggressive", common, mrt, sha1)

		if err != nil {
			return nil, err
		}

		next, err := repo.git("write-tree")

		if err != nil {
			fmt.Println("Simple merge did not work, trying automatic merge.")
			_, err = repo.git("merge-index", "-o", "git-merge-one-file", "-a")

			if err != nil {
				fmt.Println("Automated merge did not work.")
				fmt.Println("Should not be doing an Octopus.")
				return nil, errors.New("")
			}

			next, _ = repo.git("write-tree")
		}

		mrc = append(mrc, sha1)
		mrt = next
	}

	return mrc, nil
}

func octopusCommitMessage(remotes map[string]string) string {
	return "octopus commit"
}
