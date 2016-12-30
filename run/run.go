package run

import (
	"bytes"
	"errors"
	"github.com/lesfurets/git-octopus/config"
	"github.com/lesfurets/git-octopus/git"
	"github.com/lesfurets/git-octopus/test"
	"log"
	"strings"
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

	status, _ := context.Repo.Git("status", "--porcelain")

	// This is not formally required but it would be an ambiguous behaviour to let git-octopus run on unclean state.
	if len(status) != 0 {
		return errors.New("The repository has to be clean.")
	}

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

// The logic of this function is copied directly from git-merge-octopus.sh
func mergeHeads(context *OctopusContext, remotes map[string]string) ([]string, error) {
	head, _ := context.Repo.Git("rev-parse", "--verify", "-q", "HEAD")

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
