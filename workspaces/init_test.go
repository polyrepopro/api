package workspaces

import (
	"os"
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

func (s *TestSuite) Test1InitFromRemoteURL() {
	err := Init(InitArgs{
		Path: "~/.polyrepo.yaml",
		URL:  "https://raw.githubusercontent.com/polyrepopro/workspace/main/.polyrepo.yaml",
	})
	assert.NoError(s.T(), err)

	_, err = config.GetAbsoluteConfig("~/.polyrepo.yaml")
	assert.NoError(s.T(), err)

	err = os.Remove("~/.polyrepo.yaml")
	assert.NoError(s.T(), err)
}

func (s *TestSuite) Test2InitHomeDirDefault() {
	err := Init(InitArgs{
		Path: "~/.polyrepo.yaml",
	})
	assert.NoError(s.T(), err)

	_, err = config.GetAbsoluteConfig("~/.polyrepo.yaml")
	assert.NoError(s.T(), err)
}

func (s *TestSuite) Test2InitLocalDirDefault() {
	err := Init(InitArgs{
		Path: "./temp/.polyrepo.yaml",
	})
	assert.NoError(s.T(), err)

	_, err = config.GetAbsoluteConfig("./temp/.polyrepo.yaml")
	assert.NoError(s.T(), err)

	err = os.RemoveAll("./temp")
	assert.NoError(s.T(), err)
}
