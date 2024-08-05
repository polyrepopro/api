package workspaces

import (
	"log"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type CommitSuite struct {
	suite.Suite
	cfg *config.Config
}

func TestCommit(t *testing.T) {
	suite.Run(t, new(CommitSuite))
}

func (s *CommitSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = config.GetAbsoluteConfig("~/.polyrepo.test.yaml")
	if err != nil {
		log.Fatalf("failed to get absolute config: %v", err)
	}
}

func (s *CommitSuite) Test1Commit() {
	res, errs := Commit(CommitArgs{
		Workspace: &(*s.cfg.Workspaces)[0],
	})
	if len(errs) > 0 {
		log.Printf("errs: %v", errs)
	}
	assert.Equal(s.T(), 0, len(errs))
	assert.Equal(s.T(), 1, len(res))
}
