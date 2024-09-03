package config

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	path string
}

func (s *TestSuite) SetupTest() {
	test.Setup()
	s.path = "~/.polyrepo.yaml"
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TearDownSuite() {

}

func (s *TestSuite) Test1GetAbsoluteConfig() {
	_, err := GetAbsoluteConfig(s.path)
	assert.NoError(s.T(), err)
}

func (s *TestSuite) Test2GetRelativeConfig() {
	config, err := GetRelativeConfig()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), files.ExpandPath(s.path), config.Path)
}

func (s *TestSuite) Test3GetConfigWithPath() {
	config, err := GetConfig(s.path)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), files.ExpandPath(s.path), config.Path)
}

func (s *TestSuite) Test4GetConfigWithoutPath() {
	config, err := GetConfig("")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), files.ExpandPath(s.path), config.Path)
}
