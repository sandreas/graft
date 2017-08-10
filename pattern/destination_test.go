package pattern

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewDestinationPattern(t *testing.T) {
	expect := assert.New(t)

	sourcePattern := NewDestinationPattern("../data/tmp/new$1_file")
	expect.Equal("../data/tmp", sourcePattern.Path)
	expect.Equal("new${1}_file", sourcePattern.Pattern)
}