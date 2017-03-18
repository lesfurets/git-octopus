package config

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/lesfurets/git-octopus/git"
)

type OctopusConfig struct {
	PrintVersion     bool
	DoCommit         bool
	RecursiveMode    bool
	ChunkSize        int
	ExcludedPatterns []string
	Patterns         []string
}

type excluded_patterns []string

func (e *excluded_patterns) String() string {
	return strings.Join(*e, ",")
}

func (e *excluded_patterns) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func GetOctopusConfig(repo *git.Repository, args []string) (*OctopusConfig, error) {

	var printVersion, noCommitArg, commitArg, recursiveArg bool
	var chunkSizeArg int
	var excludedPatternsArg excluded_patterns

	var commandLine = flag.NewFlagSet("git-octopus", flag.ExitOnError)
	commandLine.BoolVar(&printVersion, "v", false, "prints the version of git-octopus.")
	commandLine.BoolVar(&noCommitArg, "n", false, "leaves the repository back to HEAD.")
	commandLine.BoolVar(&commitArg, "c", false, "Commit the resulting merge in the current branch.")
	commandLine.BoolVar(&recursiveArg, "r", false, "merge using a traditional recursive merge (implies -s 1)")
	commandLine.IntVar(&chunkSizeArg, "s", 0, "do the octopus by chunk of n branches.")
	commandLine.Var(&excludedPatternsArg, "e", "exclude branches matching the pattern.")

	commandLine.Parse(args)

	var configCommit bool

	rawConfigCommit, err := repo.Git("config", "octopus.commit")

	if err != nil {
		configCommit = true
	} else {
		configCommit, err = strconv.ParseBool(rawConfigCommit)
		if err != nil {
			return nil, errors.New("Config octopus.commit should be boolean. Given \"" + rawConfigCommit + "\"")
		}
	}

	if commitArg {
		configCommit = true
	}

	if noCommitArg {
		configCommit = false
	}

	configExcludedPatterns, _ := repo.Git("config", "--get-all", "octopus.excludePattern")

	var excludedPatterns []string

	if len(configExcludedPatterns) > 0 {
		excludedPatterns = strings.Split(configExcludedPatterns, "\n")
	}

	if len(excludedPatternsArg) > 0 {
		excludedPatterns = excludedPatternsArg
	}

	configPatterns, _ := repo.Git("config", "--get-all", "octopus.pattern")

	var patterns []string

	if len(configPatterns) > 0 {
		patterns = strings.Split(configPatterns, "\n")
	}

	if commandLine.NArg() > 0 {
		patterns = commandLine.Args()
	}

	return &OctopusConfig{
		PrintVersion:     printVersion,
		DoCommit:         configCommit,
		RecursiveMode:    recursiveArg,
		ChunkSize:        chunkSizeArg,
		ExcludedPatterns: excludedPatterns,
		Patterns:         patterns,
	}, nil
}
