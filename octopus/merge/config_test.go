package merge

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoCommit(t *testing.T) {
	//given
	args := []string{"a", "b", "c"}

	//when
	cfg, _ := GetConfig(args)

	//then
	assert.Equal(t, cfg.DoCommit, true)
}

func TestChunkMode(t *testing.T) {
	//given
	args := []string{"a", "b", "c"}

	//when
	cfg, _ := GetConfig(args)

	//then
	assert.Equal(t, cfg.ChunkSize, 0)
}

func TestExcludedPatterns(t *testing.T) {
	//given
	args := []string{"a", "b", "c"}

	//when
	cfg, _ := GetConfig(args)

	//then
	assert.Empty(t, cfg.ExcludedPatterns)
}

func TestPatterns(t *testing.T) {
	//given
	args := []string{"a", "b", "c"}

	//when
	cfg, _ := GetConfig(args)

	//then
	assert.Equal(t, cfg.Patterns, args)
}
