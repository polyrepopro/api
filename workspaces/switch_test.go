package workspaces

import (
	"log"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type SwitchSuite struct {
	suite.Suite
	cfg *config.Config
}

func TestSwitch(t *testing.T) {
	suite.Run(t, new(SwitchSuite))
}

func (s *SwitchSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = config.GetAbsoluteConfig("~/.polyrepo.yaml")
	if err != nil {
		log.Fatalf("failed to get absolute config: %v", err)
	}
}

func (s *SwitchSuite) Test1Switch() {
	errs := Switch(SwitchArgs{
		Workspace: &(*s.cfg.Workspaces)[0],
		Branch:    "madin",
	})
	if len(errs) > 0 {
		log.Printf("errs: %v", errs)
	}
	assert.Equal(s.T(), 0, len(errs))
}
