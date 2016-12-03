package main

import (
	"os/exec"
	"strings"
)

func git(repoPath string, args ...string) string {
	out, _ := exec.Command("git", append([]string{"-C", repoPath}, args...)...).Output()

	return strings.TrimSpace(string(out[:]))
}
