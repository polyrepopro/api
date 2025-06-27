package repositories

import (
	"log"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type PullSuite struct {
	suite.Suite
	cfg       *config.Config
	workspace *config.Workspace
	repo      *config.Repository
}

func (s *PullSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = config.GetAbsoluteConfig("~/workspace/cmskit.orig/workspace/.polyrepo.yaml")
	if err != nil {
		log.Fatalf("failed to get absolute config: %v", err)
	}

	s.workspace = &(*s.cfg.Workspaces)[0]
	s.repo = &(*s.workspace.Repositories)[0]
}

func TestPullSuite(t *testing.T) {
	suite.Run(t, new(PullSuite))
}

func (s *PullSuite) Test1Pull() {
	err := Pull(PullArgs{
		Workspace:  s.workspace,
		Repository: s.repo,
	})
	assert.NoError(s.T(), err)
}
