package workspaces

import (
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	cfg *config.Config
}

func (s *TestSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = Init(InitArgs{
		Path: "~/.polyrepo.yaml",
		URL:  "https://raw.githubusercontent.com/polyrepopro/workspace/main/.polyrepo.yaml",
	})
	assert.NoError(s.T(), err)
}

func TestInitSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) Test1InitFromRemoteURL() {

	assert.Equal(s.T(), s.cfg.Path, files.ExpandPath("~/.polyrepo.yaml"))

	_, err := config.GetAbsoluteConfig(s.cfg.Path)
	assert.NoError(s.T(), err)

	err = os.Remove(s.cfg.Path)
	assert.NoError(s.T(), err)
}

func (s *TestSuite) Test2InitHomeDirDefault() {
	_, err := config.GetAbsoluteConfig(s.cfg.Path)
	assert.NoError(s.T(), err)

	err = os.Remove(s.cfg.Path)
	assert.NoError(s.T(), err)
}

func (s *TestSuite) Test2InitLocalDirDefault() {
	_, err := config.GetAbsoluteConfig(s.cfg.Path)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), s.cfg)

	_, err = config.GetAbsoluteConfig(s.cfg.Path)
	assert.NoError(s.T(), err)

	err = os.RemoveAll("./temp")
	assert.NoError(s.T(), err)
}
