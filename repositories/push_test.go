package repositories

import (
	"log"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type PushSuite struct {
	suite.Suite
	cfg       *config.Config
	workspace *config.Workspace
	repo      *config.Repository
}

func (s *PushSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = config.GetAbsoluteConfig("~/.polyrepo.test.yaml")
	if err != nil {
		log.Fatalf("failed to get absolute config: %v", err)
	}

	s.workspace = &(*s.cfg.Workspaces)[0]
	s.repo = &(*s.workspace.Repositories)[0]
}

func TestPushSuite(t *testing.T) {
	suite.Run(t, new(PushSuite))
}

func (s *PushSuite) Test1Push() {
	err := Push(PushArgs{
		Workspace:  s.workspace,
		Repository: s.repo,
	})
	assert.NoError(s.T(), err)
}
