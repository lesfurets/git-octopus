package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
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

func (repo *repository) git(args ...string) string {
	return git(repo.path, args...)
}

func (repo *repository) writeFile(name string, lines ...string) {
	fileName := filepath.Join(repo.path, name)
	_, err := os.Stat(fileName)

	var file *os.File
	if os.IsNotExist(err) {
		file, _ = os.Create(fileName)
	} else {
		file, _ = os.Open(fileName)
	}

	file.WriteString(strings.Join(lines, "\n"))
}
