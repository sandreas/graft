package pattern

import (
	"path/filepath"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/spf13/afero"
)

var mockFs afero.Fs
func init() {
	mockFs = afero.NewMemMapFs()
}

func TestAbsoluteWindows(t *testing.T) {

	expect := assert.New(t)

	abs := "C:/"
	absWithPattern := abs + "NotExisting/*.txt"
	sourcePattern := NewBasePattern(mockFs, absWithPattern)
	expect.Equal(filepath.ToSlash(abs), sourcePattern.Path)
	expect.Equal("NotExisting/*.txt", sourcePattern.Pattern)

}
