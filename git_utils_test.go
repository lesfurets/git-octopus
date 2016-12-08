package main

import (
	"reflect"
	"testing"
)

func TestParseLsRemoteEmpty(t *testing.T) {
	if !reflect.DeepEqual(parseLsRemote(""), map[string]string{}) {
		t.Error("Excpected to be non nil")
	}
}

func TestParseLsRemote(t *testing.T) {
	lsRemoteOutput := "d8dd4eadaf3c1075eff3b7d4fe6bec5fbfe76b4c	refs/heads/master\n" +
		"5b2b1bf1cdf1150f34bd5809a038b292dc560998	refs/heads/go_rewrite"
	if !reflect.DeepEqual(
		parseLsRemote(lsRemoteOutput),
		map[string]string{
			"refs/heads/master":     "d8dd4eadaf3c1075eff3b7d4fe6bec5fbfe76b4c",
			"refs/heads/go_rewrite": "5b2b1bf1cdf1150f34bd5809a038b292dc560998"}) {
		t.Error()
	}
}
