package pattern_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/testhelpers"
	"github.com/sandreas/graft/pattern"
	"os"
)

func TestNewDestinationPattern(t *testing.T) {
	expect := assert.New(t)

	mockFs := testhelpers.MockFileSystem(map[string]string{
		"data/tmp/":         "",
	})
	sep := string(os.PathSeparator)
	destinationPattern := pattern.NewDestinationPattern(mockFs, "data/tmp/new$1_file")
	expect.Equal("data" +sep+"tmp", destinationPattern.Path)
	expect.Equal("new${1}_file", destinationPattern.Pattern)
}