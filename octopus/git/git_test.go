package git

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestParseLsRemoteEmpty(t *testing.T) {
	assert.NotNil(t, ParseLsRemote(""))
	assert.Equal(t, []LsRemoteEntry{}, ParseLsRemote(""))
}

func TestParseLsRemote(t *testing.T) {
	lsRemoteOutput := "d8dd4eadaf3c1075eff3b7d4fe6bec5fbfe76b4c	refs/heads/master\n" +
		"5b2b1bf1cdf1150f34bd5809a038b292dc560998	refs/heads/go_rewrite"
	expected := []LsRemoteEntry{
		{Ref: "refs/heads/master", Sha1: "d8dd4eadaf3c1075eff3b7d4fe6bec5fbfe76b4c"},
		{Ref: "refs/heads/go_rewrite", Sha1: "5b2b1bf1cdf1150f34bd5809a038b292dc560998"},
	}
	assert.Equal(t, expected, ParseLsRemote(lsRemoteOutput))
}

func TestGitCommand(t *testing.T) {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")
	defer os.RemoveAll(dir)

	repo := Repository{Path: dir}

	repo.Git("init")

	_, err := os.Stat(filepath.Join(dir, ".git"))

	assert.Nil(t, err)
}

func TestGitError(t *testing.T) {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")
	defer os.RemoveAll(dir)

	repo := Repository{Path: dir}

	repo.Git("init")

	_, err := repo.Git("rev-parse", "HEAD")

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(),
			"ambiguous argument 'HEAD': unknown revision or path not in the working tree.")
	}
}
