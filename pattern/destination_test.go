package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/testhelpers"
	"github.com/sandreas/graft/pattern"
)

func TestNewDestinationPattern(t *testing.T) {
	expect := assert.New(t)

	mockFs := testhelpers.MockFileSystem(map[string]string{
		"data/tmp/":         "",
	})

	destinationPattern := pattern.NewDestinationPattern(mockFs, "data/tmp/new$1_file")
	expect.Equal("data/tmp", destinationPattern.Path)
	expect.Equal("new${1}_file", destinationPattern.Pattern)
}