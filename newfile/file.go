package newfile

import (
	"github.com/sandreas/graft/newpattern"
	"os"
	"path/filepath"
	"github.com/sandreas/graft/newmatcher"
	"github.com/sandreas/graft/newprogress"
)


func FindFilesBySourcePattern(p newpattern.SourcePattern, matcher newmatcher.MatcherInterface, progressHandler newprogress.ProgressHandler) ([]string, error) {
	var m []string
	if p.IsFile() {
		m = append(m, p.Path)

		progressHandler.IncreaseMatches()
		progressHandler.Finish()
		return m, nil
	}

	filepath.Walk(p.Path, func(innerPath string, info os.FileInfo, err error) error {
		if innerPath == "." || innerPath == ".." {
			return nil
		}


		normalizedInnerPath := filepath.ToSlash(innerPath)
		if info.IsDir() {
			normalizedInnerPath += "/"
		}

		if matcher.Matches(normalizedInnerPath) {
			m = append(m, normalizedInnerPath)
			progressHandler.IncreaseMatches()
		} else {
			progressHandler.IncreaseItems()
		}

		return nil
	})

	progressHandler.Finish()

	return m, nil
}

