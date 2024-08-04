package repositories

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type AddSuite struct {
	suite.Suite
	workspace *config.Workspace
}

func (s *AddSuite) SetupTest() {
	test.Setup()

	var err error

	config, err := config.GetRelativeConfig()
	assert.NoError(s.T(), err)

	s.workspace, err = config.GetWorkspaceByWorkingDir()
	assert.NoError(s.T(), err)
}

func TestAddSuite(t *testing.T) {
	suite.Run(t, new(AddSuite))
}

func (s *AddSuite) TearDownSuite() {
	err := Remove(config.Repository{
		Path: "test/test",
	})
	assert.NoError(s.T(), err)
}

func (s *AddSuite) Test1Add() {
	err := Add(config.Repository{
		Branch: "main",
		Path:   "test/test",
		URL:    "https://github.com/risersh/example-test-repo.git",
	})
	assert.NoError(s.T(), err)
}

func (s *AddSuite) Test2Get() {
	repository, err := Get("test/test")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test/test", repository.Path)
}
