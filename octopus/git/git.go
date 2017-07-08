package git

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type LsRemoteEntry struct {
	Ref  string
	Sha1 string
}

// Takes the output of git-ls-remote. Returns a map refsname => sha1
func ParseLsRemote(lsRemoteOutput string) []LsRemoteEntry {
	result := []LsRemoteEntry{}

	if len(lsRemoteOutput) == 0 {
		return result
	}

	scanner := bufio.NewScanner(strings.NewReader(lsRemoteOutput))

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "\t")
		result = append(result, LsRemoteEntry{Ref: split[1], Sha1: split[0]})
	}

	return result
}

type Repository struct {
	Path string
}

func (repo *Repository) Git(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", repo.Path}, args...)...)
	errOut := bytes.NewBufferString("")
	cmd.Stderr = errOut
	out, err := cmd.Output()

	if err != nil {
		return "", errors.New(errOut.String())
	}

	stringOut := strings.TrimSpace(string(out[:]))

	return stringOut, err
}
