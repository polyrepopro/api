package repositories

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	workspace *config.Workspace
}

func (s *TestSuite) SetupTest() {
	test.Setup()

	var err error

	config, err := config.GetRelativeConfig()
	assert.NoError(s.T(), err)

	s.workspace, err = config.GetWorkspaceByWorkingDir()
	assert.NoError(s.T(), err)
}

func TestDoctorSuiteRun(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TearDownSuite() {

}

func (s *TestSuite) TestDoctor() {
	for _, repository := range *s.workspace.Repositories {
		err := Doctor(&repository, s.workspace)
		assert.NoError(s.T(), err)
	}
}
