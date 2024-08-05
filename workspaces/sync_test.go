package workspaces

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type SyncSuite struct {
	suite.Suite
}

func (s *SyncSuite) SetupTest() {
	test.Setup()
}

// func (s *SyncSuite) TearDownTest() {
// 	err := os.Remove("~/.polyrepo.yaml")
// 	assert.NoError(s.T(), err)
// }

func TestSync(t *testing.T) {
	suite.Run(t, new(SyncSuite))
}

func (s *SyncSuite) Test1Sync() {

	err := Init(InitArgs{
		Path: "~/.polyrepo.yaml",
		URL:  "https://raw.githubusercontent.com/polyrepopro/workspace/main/.polyrepo.yaml",
	})
	assert.NoError(s.T(), err)

	err = Sync(SyncArgs{})
	assert.NoError(s.T(), err)
}
