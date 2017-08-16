package testhelpers

import (
	"github.com/spf13/afero"
	"strings"
)

func MockFileSystem(initialFiles map[string]string) afero.Fs {
	mockFs := afero.NewMemMapFs()

	for key, value := range initialFiles {
		if strings.HasSuffix(key, "/") || strings.HasSuffix(key, "\\") {
			mockFs.Mkdir(key, 0644)
		} else {
			afero.WriteFile(mockFs,key, []byte(value), 0755)
		}
	}
	return mockFs
}