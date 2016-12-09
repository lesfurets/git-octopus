package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoCommit(t *testing.T) {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	// GIVEN no config, no option
	// WHEN
	octopusConfig, _ := getOctopusConfig(repo, nil)

	// THEN doCommit should be true
	assert.True(t, octopusConfig.doCommit)

	// GIVEN config to false, no option
	repo.git("config", "octopus.commit", "false")
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, nil)

	// THEN doCommit should be false
	assert.False(t, octopusConfig.doCommit)

	// GIVEN config to false, -c option true
	repo.git("config", "octopus.commit", "false")
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, []string{"-c"})

	// THEN  doCommit should be true
	assert.True(t, octopusConfig.doCommit)

	// GIVEN config to true, -n option true
	repo.git("config", "octopus.commit", "true")
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, []string{"-n"})

	// THEN  doCommit should be false
	assert.False(t, octopusConfig.doCommit)
}

func TestChunkMode(t *testing.T) {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	// GIVEN No option
	// WHEN
	octopusConfig, _ := getOctopusConfig(repo, nil)

	// THEN chunkSize should be 0
	assert.Equal(t, 0, octopusConfig.chunkSize)

	// GIVEN option -s 5
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, []string{"-s", "5"})

	// THEN chunkSize should be 5
	assert.Equal(t, 5, octopusConfig.chunkSize)
}

func TestExcludedPatterns(t *testing.T) {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	// GIVEN no config, no option
	// WHEN
	octopusConfig, _ := getOctopusConfig(repo, nil)

	// THEN excludedPatterns should be empty
	assert.Empty(t, octopusConfig.excludedPatterns)

	// GIVEN excludePattern config, no option
	repo.git("config", "octopus.excludePattern", "excluded/*")
	repo.git("config", "--add", "octopus.excludePattern", "excluded_branch")
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, nil)

	// THEN excludedPatterns should be set
	assert.Equal(t, []string{"excluded/*", "excluded_branch"}, octopusConfig.excludedPatterns)

	// GIVEN excludePattern config (from previous assertion), option given
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, []string{"-e", "override_excluded"})

	// THEN option should take precedence
	assert.Equal(t, []string{"override_excluded"}, octopusConfig.excludedPatterns)
}

func TestPatterns(t *testing.T) {
	repo := createTestRepo()
	defer cleanupTestRepo(repo)

	// GIVEN no config, no option
	// WHEN
	octopusConfig, _ := getOctopusConfig(repo, nil)

	// THEN excludedPatterns should be empty
	assert.Empty(t, octopusConfig.patterns)

	// GIVEN config, no argument.
	repo.git("config", "octopus.pattern", "test")
	repo.git("config", "--add", "octopus.pattern", "test2")
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, nil)

	// THEN patterns should be set
	assert.Equal(t, []string{"test", "test2"}, octopusConfig.patterns)

	// GIVEN config (from previous assertion), argument given
	// WHEN
	octopusConfig, _ = getOctopusConfig(repo, []string{"arg1", "arg2"})

	// THEN arguments should take precedence
	assert.Equal(t, []string{"arg1", "arg2"}, octopusConfig.patterns)
}
