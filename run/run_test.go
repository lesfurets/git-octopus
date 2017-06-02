package run

import (
	"github.com/lesfurets/git-octopus/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeOneParent(t *testing.T) {
	//given
	context, _ := CreateTestContext()
	repo := context.Repo
	defer test.Cleanup(repo)

	writeFile(repo, "test", "")
	repo.Git("add", "test")
	repo.Git("commit", "-a -m \"commit test\"")
	repo.Git("checkout", "-b", "qa")
	repo.Git("checkout", "-b", "feature/test")
	writeFile(repo, "testFeature", "")
	repo.Git("add", "testFeature")
	repo.Git("commit", "-a -m \"add testFeature\"")
	testFeatureSha1, _ := repo.Git("rev-parse", "HEAD")
	repo.Git("checkout", "master")

	//when
	Run(context, "", "feature/* qa")

	//then
	actual, _ := repo.Git("branch", "--contains", testFeatureSha1)
	assert.Contains(t, actual, "feature/test")
}
