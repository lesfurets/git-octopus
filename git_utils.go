package main

import (
	"os/exec"
	"strings"
	"bufio"
	"path/filepath"
	"runtime"
	"io/ioutil"
	"fmt"
)

func git(repoPath string, args ...string) string {
	out, _ := exec.Command("git", append([]string{"-C", repoPath}, args...)...).Output()

	return strings.TrimSpace(string(out[:]))
}

// Takes the output of git-ls-remote. Returns a map sha1 => refsname
func parseLsRemote(lsRemoteOutput string) map[string]string {
	result := map[string]string{}

	if len(lsRemoteOutput) == 0 {
		return result
	}

	scanner := bufio.NewScanner(strings.NewReader(lsRemoteOutput))

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "\t")

		result[split[0]] = split[1]
	}

	return result
}


type repository struct {
	path string
}

func createRepo() *repository {
	dir, _ := ioutil.TempDir("", "git-octopus-test-")

	repository := repository{ path: dir}

	git(dir, "init")

	return &repository
}

func (repo *repository) git(args ...string) string {
	return git(repo.path, args...)
}

func (repo *repository) run(script string) {
	_, file, _, _ := runtime.Caller(0)
	scriptPath := filepath.Join(filepath.Dir(file), script)
	cmd := exec.Command(scriptPath)
	cmd.Dir = repo.path
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}