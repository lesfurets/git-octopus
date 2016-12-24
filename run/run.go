package run

import (
	"lesfurets/git-octopus/git"
	"log"
	"errors"
	"lesfurets/git-octopus/config"
	"strings"
	"bytes"
	"lesfurets/git-octopus/test"
)

type OctopusContext struct {
	Repo   *git.Repository
	Logger *log.Logger
}

func Run(context *OctopusContext, args ...string) error {

	octopusConfig, err := config.GetOctopusConfig(context.Repo, args)

	if err != nil {
		return err
	}

	if octopusConfig.PrintVersion {
		context.Logger.Println("2.0")
		return nil
	}

	if len(octopusConfig.Patterns) == 0 {
		context.Logger.Println("Nothing to merge. No pattern given")
		return nil
	}

	branchList := resolveBranchList(context.Repo, octopusConfig.Patterns, octopusConfig.ExcludedPatterns)

	if len(branchList) == 0 {
		context.Logger.Printf("No branch matching \"%v\" were found\n", strings.Join(octopusConfig.Patterns, " "))
		return nil
	}

	initialHeadCommit, _ := context.Repo.Git("rev-parse", "HEAD")

	parents, err := mergeHeads(context, branchList)

	if !octopusConfig.DoCommit {
		context.Repo.Git("reset", "-q", "--hard", initialHeadCommit)
	}

	if err != nil {
		return err
	}

	if octopusConfig.DoCommit {
		tree, _ := context.Repo.Git("write-tree")
		commit, _ := context.Repo.Git("commit-tree", "-p", strings.Join(parents, " -p "), "-m", octopusCommitMessage(branchList), tree)
		context.Repo.Git("update-ref", "HEAD", commit)
	}

	return nil
}

// The logic of this function is copied directly from git-merge-octopus.sh
func mergeHeads(context *OctopusContext, remotes map[string]string) ([]string, error) {
	head, _ := context.Repo.Git("rev-parse", "--verify", "-q", "HEAD")

	alreadyUpToDate := true
	for _, sha1 := range remotes {
		_, err := context.Repo.Git("merge-base", "--is-ancestor", sha1, "HEAD")
		if err != nil {
			alreadyUpToDate = false
		}
	}
	// This prevents git-octopus to create a commit when there's nothing to merge,
	// i.e. no feature branches but only master.
	if alreadyUpToDate {
		context.Logger.Println("Already up to date")
		return nil, nil
	}

	mrc := []string{head}
	mrt, _ := context.Repo.Git("write-tree")
	nonFfMerge := false

	for prettyRemoteName, sha1 := range remotes {

		common, err := context.Repo.Git(append([]string{"merge-base", "--all", sha1}, mrc...)...)

		if err != nil {
			return nil, errors.New("Unable to find common commit with " + prettyRemoteName)
		}

		if common == sha1 {
			context.Logger.Println("Already up-to-date with " + prettyRemoteName)
			continue
		}

		if len(mrc) == 1 && common == mrc[0] && !nonFfMerge {
			context.Logger.Println("Fast-forwarding to: " + prettyRemoteName)
			_, err := context.Repo.Git("read-tree", "-u", "-m", head, sha1)

			if err != nil {
				return nil, nil
			}

			mrc[0] = sha1
			mrt, _ = context.Repo.Git("write-tree")
			continue
		}

		nonFfMerge = true

		context.Logger.Println("Trying simple merge with " + prettyRemoteName)

		_, err = context.Repo.Git("read-tree", "-u", "-m", "--aggressive", common, mrt, sha1)

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

		mrc = append(mrc, sha1)
		mrt = next
	}

	return mrc, nil
}

func octopusCommitMessage(remotes map[string]string) string {
	return "octopus commit"
}


func CreateTestContext() (*OctopusContext, *bytes.Buffer) {
	dir := test.CreateTempDir()

	repo := git.Repository{Path: dir}

	repo.Git("init")
	repo.Git("commit", "--allow-empty", "-m\"first commit\"")

	out := bytes.NewBufferString("")

	context := OctopusContext{
		Repo:   &repo,
		Logger: log.New(out, "", 0),
	}

	return &context, out
}
