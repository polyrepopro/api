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

type SyncSuite struct {
	suite.Suite
	cfg *config.Config
}

func TestSync(t *testing.T) {
	suite.Run(t, new(SyncSuite))
}

func (s *SyncSuite) SetupTest() {
	test.Setup()

	var err error
	s.cfg, err = Init(InitArgs{
		Path: "~/.polyrepo.yaml",
		URL:  "https://raw.githubusercontent.com/polyrepopro/workspace/main/.polyrepo.yaml",
	})
	assert.NoError(s.T(), err)
}

func (s *SyncSuite) TearDownTest() {
	err := os.Remove(s.cfg.Path)
	assert.NoError(s.T(), err)

	for _, workspace := range *s.cfg.Workspaces {
		err = os.RemoveAll(files.ExpandPath(workspace.Path))
		assert.NoError(s.T(), err)
	}
}

func (s *SyncSuite) Test1Sync() {

	assert.NotNil(s.T(), s.cfg)

	_, errs := SyncAll(nil)
	assert.Equal(s.T(), 0, len(errs))
}
