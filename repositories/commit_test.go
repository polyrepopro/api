package repositories

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/polyrepopro/api/util"
	"github.com/stretchr/testify/suite"
)

type CommitSuite struct {
	suite.Suite
	cfg       *config.Config
	workspace *config.Workspace
	repo      *config.Repository
}

func (s *CommitSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = config.GetAbsoluteConfig("~/.polyrepo.test.yaml")
	if err != nil {
		log.Fatalf("failed to get absolute config: %v", err)
	}

	s.workspace = &(*s.cfg.Workspaces)[0]
	s.repo = &(*s.workspace.Repositories)[0]
}

func TestCommitSuite(t *testing.T) {
	suite.Run(t, new(CommitSuite))
}

func (s *CommitSuite) Test1Commit() {
	testFilePath := fmt.Sprintf("%s/%s/test_commits.txt", util.ExpandPath(s.workspace.Path), s.repo.Path)
	testContent := fmt.Sprintf("Test commit @ %s", time.Now())

	err := os.WriteFile(testFilePath, []byte(testContent), 0644)
	assert.NoError(s.T(), err)

	hash, err := Commit(CommitArgs{
		Workspace:  s.workspace,
		Repository: s.repo,
		Message:    fmt.Sprintf("Test commit @ %s", time.Now()),
	})
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), hash)
}
