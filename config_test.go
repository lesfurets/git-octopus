package main

import (
	"reflect"
	"testing"
)

func TestDoCommit(t *testing.T) {
	repo := createRepo()

	// No config, no option. doCommit should be true
	octopusConfig := getOctopusConfig(repo.path, nil)

	if !octopusConfig.doCommit {
		t.Error("Expected doCommit = true, got", octopusConfig.doCommit)
	}

	// Config to false, no option. doCommit should be false
	repo.git("config", "octopus.commit", "false")
	octopusConfig = getOctopusConfig(repo.path, nil)

	if octopusConfig.doCommit {
		t.Error("Expected doCommit = false, got", octopusConfig.doCommit)
	}

	// Config to false, -c option takes precedence. doCommit should be true
	repo.git("config", "octopus.commit", "false")
	octopusConfig = getOctopusConfig(repo.path, []string{"-c"})

	if !octopusConfig.doCommit {
		t.Error("Expected doCommit = true, got", octopusConfig.doCommit)
	}

	// Config to true, -n option takes precedence. doCommit should be false
	repo.git("config", "octopus.commit", "true")
	octopusConfig = getOctopusConfig(repo.path, []string{"-n"})

	if octopusConfig.doCommit {
		t.Error("Expected doCommit = false, got", octopusConfig.doCommit)
	}
}

func TestChunkMode(t *testing.T) {
	repo := createRepo()

	// No option. chunkSize should be set to 0
	octopusConfig := getOctopusConfig(repo.path, nil)

	if octopusConfig.chunkSize != 0 {
		t.Error("Expected chunkSize = 0, got", octopusConfig.chunkSize)
	}

	// -s 5, chunkSize should be set to 5
	octopusConfig = getOctopusConfig(repo.path, []string{"-s", "5"})

	if octopusConfig.chunkSize != 5 {
		t.Error("Expected chunkSize = 5, got", octopusConfig.chunkSize)
	}
}

func TestExcludedPatterns(t *testing.T) {
	repo := createRepo()

	// No config, no option. excludedPatterns should be empty
	octopusConfig := getOctopusConfig(repo.path, nil)

	if len(octopusConfig.excludedPatterns) > 0 {
		t.Error("Expected excludedPatterns to be empty")
	}

	// Config given, no option. excludedPatterns should be set
	repo.git("config", "octopus.excludePattern", "excluded/*")
	repo.git("config", "--add", "octopus.excludePattern", "excluded_branch")
	octopusConfig = getOctopusConfig(repo.path, nil)

	if !reflect.DeepEqual(octopusConfig.excludedPatterns, []string{"excluded/*", "excluded_branch"}) {
		t.Error("actual excludedPatterns:", octopusConfig.excludedPatterns)
	}

	// Config given (from previous assertion), option given. Option should take precedence
	octopusConfig = getOctopusConfig(repo.path, []string{"-e", "override_excluded"})

	if !reflect.DeepEqual(octopusConfig.excludedPatterns, []string{"override_excluded"}) {
		t.Error("excludedPatterns", octopusConfig.excludedPatterns)
	}
}

func TestPatterns(t *testing.T) {
	repo := createRepo()

	// No config, no option. excludedPatterns should be empty
	octopusConfig := getOctopusConfig(repo.path, nil)

	if len(octopusConfig.patterns) > 0 {
		t.Error("Expected patterns to be empty")
	}

	// Config given, no argument. patterns should be set
	repo.git("config", "octopus.pattern", "test")
	repo.git("config", "--add", "octopus.pattern", "test2")
	octopusConfig = getOctopusConfig(repo.path, nil)

	if !reflect.DeepEqual(octopusConfig.patterns, []string{"test", "test2"}) {
		t.Error("actual patterns:", octopusConfig.patterns)
	}

	// Config given (from previous assertion), argument given. Arguments should take precedence
	octopusConfig = getOctopusConfig(repo.path, []string{"arg1", "arg2"})

	if !reflect.DeepEqual(octopusConfig.patterns, []string{"arg1", "arg2"}) {
		t.Error("actual patterns:", octopusConfig.patterns)
	}
}
