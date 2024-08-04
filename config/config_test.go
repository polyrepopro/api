package config

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/test"
	"github.com/polyrepopro/api/util"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	path string
}

func (s *TestSuite) SetupTest() {
	test.Setup()
	s.path = "~/workspace/.polyrepo.yaml"
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
	assert.Equal(s.T(), util.ExpandPath(s.path), config.Path)
}
