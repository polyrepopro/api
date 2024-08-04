package workspaces

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (s *TestSuite) SetupTest() {
	test.Setup()
}

func TestInitSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TearDownSuite() {

}

func (s *TestSuite) TestInit() {
	err := Init(InitArgs{
		Path: "~/workspace/.polyrepo.yaml",
		URL:  "https://raw.githubusercontent.com/polyrepopro/workspace/main/.poly.yaml",
	})
	assert.NoError(s.T(), err)

	_, err = config.GetAbsoluteConfig("~/workspace/.polyrepo.yaml")
	assert.NoError(s.T(), err)
}
