package workspaces

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
	cfg *config.Config
}

func TestPull(t *testing.T) {
	suite.Run(t, new(PullSuite))
}

func (s *PullSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = config.GetAbsoluteConfig("~/workspace/cmskit.orig/workspace/.polyrepo.yaml")
	if err != nil {
		log.Fatalf("failed to get absolute config: %v", err)
	}
}

func (s *PullSuite) Test1Pull() {
	errs := Pull(PullArgs{
		Workspace: &(*s.cfg.Workspaces)[0],
	})
	if len(errs) > 0 {
		log.Printf("errs: %v", errs)
	}
	assert.Equal(s.T(), 0, len(errs))
}
