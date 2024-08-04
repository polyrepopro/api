package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	path := "~/test"
	expanded, err := ExpandPath(path)
	assert.NoError(t, err)
	assert.Equal(t, "/Users/matthewdavis/test", expanded)
}
