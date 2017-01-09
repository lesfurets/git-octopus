package run

import (
	"bytes"
	"github.com/lesfurets/git-octopus/test"
	"github.com/lesfurets/git-octopus/git"
	"fmt"
	"log"
)

func CreateTestContext() (*OctopusContext, *bytes.Buffer) {
	dir := test.CreateTempDir()

	repo := git.Repository{Path: dir}

	repo.Git("init")
	repo.Git("config", "user.name", "gotest")
	repo.Git("config", "user.email", "gotest@golang.com")
	_, err := repo.Git("commit", "--allow-empty", "-m\"first commit\"")

	if err != nil {
		fmt.Println("There's something wrong with the git installation:")
		fmt.Println(err.Error())
	}

	out := bytes.NewBufferString("")

	context := OctopusContext{
		Repo:   &repo,
		Logger: log.New(out, "", 0),
	}

	return &context, out
}
