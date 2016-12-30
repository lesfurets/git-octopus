package git

import (
	"github.com/stretchr/testify/assert"
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
