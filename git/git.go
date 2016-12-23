package git

import (
	"bufio"
	"os/exec"
	"strings"
)

// Takes the output of git-ls-remote. Returns a map refsname => sha1
func ParseLsRemote(lsRemoteOutput string) map[string]string {
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

type Repository struct {
	Path string
}

func (repo *Repository) Git(args ...string) (string, error) {
	out, err := exec.Command("git", append([]string{"-C", repo.Path}, args...)...).Output()

	stringOut := strings.TrimSpace(string(out[:]))

	return stringOut, err
}
