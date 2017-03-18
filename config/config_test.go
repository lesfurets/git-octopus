package config

import (
	"testing"

	"github.com/lesfurets/git-octopus/git"
	"github.com/lesfurets/git-octopus/test"
	"github.com/stretchr/testify/assert"
)

func createTestRepo() *git.Repository {
	dir := test.CreateTempDir()

	repo := git.Repository{Path: dir}

	repo.Git("init")

	return &repo
}

func TestDoCommit(t *testing.T) {
	repo := createTestRepo()
	defer test.Cleanup(repo)

	// GIVEN no config, no option
	// WHEN
	octopusConfig, err := GetOctopusConfig(repo, nil)

	// THEN doCommit should be true
	assert.True(t, octopusConfig.DoCommit)
	assert.Nil(t, err)

	// GIVEN config to false, no option
	repo.Git("config", "octopus.commit", "false")
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, nil)

	// THEN doCommit should be false
	assert.False(t, octopusConfig.DoCommit)
	assert.Nil(t, err)

	// Config to 0, no option. doCommit should be true
	repo.Git("config", "octopus.commit", "0")
	octopusConfig, err = GetOctopusConfig(repo, nil)

	assert.False(t, octopusConfig.DoCommit)
	assert.Nil(t, err)

	// GIVEN config to false, -c option true
	repo.Git("config", "octopus.commit", "false")
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, []string{"-c"})

	// THEN  doCommit should be true
	assert.True(t, octopusConfig.DoCommit)
	assert.Nil(t, err)

	// GIVEN config to true, -n option true
	repo.Git("config", "octopus.commit", "true")
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, []string{"-n"})

	// THEN  doCommit should be false
	assert.False(t, octopusConfig.DoCommit)
	assert.Nil(t, err)
}

func TestChunkMode(t *testing.T) {
	repo := createTestRepo()
	defer test.Cleanup(repo)

	// GIVEN No option
	// WHEN
	octopusConfig, err := GetOctopusConfig(repo, nil)

	// THEN chunkSize should be 0
	assert.Equal(t, 0, octopusConfig.ChunkSize)
	assert.Nil(t, err)

	// GIVEN option -s 5
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, []string{"-s", "5"})

	// THEN chunkSize should be 5
	assert.Equal(t, 5, octopusConfig.ChunkSize)
	assert.Nil(t, err)
}

func TestExcludedPatterns(t *testing.T) {
	repo := createTestRepo()
	defer test.Cleanup(repo)

	// GIVEN no config, no option
	// WHEN
	octopusConfig, err := GetOctopusConfig(repo, nil)

	// THEN excludedPatterns should be empty
	assert.Empty(t, octopusConfig.ExcludedPatterns)
	assert.Nil(t, err)

	// GIVEN excludePattern config, no option
	repo.Git("config", "octopus.excludePattern", "excluded/*")
	repo.Git("config", "--add", "octopus.excludePattern", "excluded_branch")
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, nil)

	// THEN excludedPatterns should be set
	assert.Equal(t, []string{"excluded/*", "excluded_branch"}, octopusConfig.ExcludedPatterns)
	assert.Nil(t, err)

	// GIVEN excludePattern config (from previous assertion), option given
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, []string{"-e", "override_excluded"})

	// THEN option should take precedence
	assert.Equal(t, []string{"override_excluded"}, octopusConfig.ExcludedPatterns)
	assert.Nil(t, err)
}

func TestPatterns(t *testing.T) {
	repo := createTestRepo()
	defer test.Cleanup(repo)

	// GIVEN no config, no option
	// WHEN
	octopusConfig, err := GetOctopusConfig(repo, nil)

	// THEN excludedPatterns should be empty
	assert.Empty(t, octopusConfig.Patterns)
	assert.Nil(t, err)

	// GIVEN config, no argument.
	repo.Git("config", "octopus.pattern", "test")
	repo.Git("config", "--add", "octopus.pattern", "test2")
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, nil)

	// THEN patterns should be set
	assert.Equal(t, []string{"test", "test2"}, octopusConfig.Patterns)
	assert.Nil(t, err)

	// GIVEN config (from previous assertion), argument given
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, []string{"arg1", "arg2"})

	// THEN arguments should take precedence
	assert.Equal(t, []string{"arg1", "arg2"}, octopusConfig.Patterns)
	assert.Nil(t, err)
}

func TestRecurisveMode(t *testing.T) {
	repo := createTestRepo()
	defer test.Cleanup(repo)

	// GIVEN No option
	// WHEN
	octopusConfig, err := GetOctopusConfig(repo, nil)

	// THEN RecursiveMode should be false
	assert.False(t, octopusConfig.RecursiveMode)
	assert.Nil(t, err)

	// GIVEN option -r
	// WHEN
	octopusConfig, err = GetOctopusConfig(repo, []string{"-r"})

	// THEN RecursiveMode should be true
	assert.True(t, octopusConfig.RecursiveMode)
	assert.Nil(t, err)
}
