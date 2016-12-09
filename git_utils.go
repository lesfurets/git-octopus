package main

import (
	"bufio"
	"io/ioutil"
	"os/exec"
	"strings"
)

func git(repoPath string, args ...string) string {
	out, _ := exec.Command("git", append([]string{"-C", repoPath}, args...)...).Output()

	return strings.TrimSpace(string(out[:]))
}

// Takes the output of git-ls-remote. Returns a map refsname => sha1
func parseLsRemote(lsRemoteOutput string) map[string]string {
	result := map[string]string{}

	if len(lsRemoteOutput) == 0 {
		return result
	}

	scanner := bufio.NewScanner(strings.NewReader(lsRemoteOutput))

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "\t")

		result[split[1]] = split[0]
	}

	return result
}

type repository struct {
	path string
}

func createRepo() *repository {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")

	repo := repository{path: dir}

	repo.git("init")

	return &repo
}

func (repo *repository) git(args ...string) string {
	return git(repo.path, args...)
}
