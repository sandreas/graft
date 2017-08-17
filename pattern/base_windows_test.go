package pattern

import (
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/testhelpers"
)
func TestAbsoluteWindows(t *testing.T) {

	expect := assert.New(t)
	mockFs := testhelpers.MockFileSystem(map[string]string{
		"C:/":         "",
	})

	abs := "C:/"
	absWithPattern := abs + "NotExisting/*.txt"
	sourcePattern := NewBasePattern(mockFs, absWithPattern)
	expect.Equal(filepath.ToSlash(abs), sourcePattern.Path)
	expect.Equal("NotExisting/*.txt", sourcePattern.Pattern)

}
