package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	path := "~/test"
	expanded := ExpandPath(path)
	assert.Equal(t, "/Users/matthewdavis/test", expanded)
}
