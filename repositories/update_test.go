package repositories

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/polyrepopro/api/test"
	"github.com/stretchr/testify/suite"
)

type UpdateSuite struct {
	suite.Suite
}

func (s *UpdateSuite) SetupTest() {
	test.Setup()
}

// func (s *SyncSuite) TearDownTest() {
// 	err := os.Remove("~/.polyrepo.yaml")
// 	assert.NoError(s.T(), err)
// }

func TestUpdate(t *testing.T) {
	suite.Run(t, new(UpdateSuite))
}

func (s *UpdateSuite) Test1Update() {

	err := Update(nil, nil)
	assert.NoError(s.T(), err)
}
