package newfile

import (
	"github.com/sandreas/graft/newpattern"
	"os"
	"path/filepath"
	"github.com/sandreas/graft/newmatcher"
)

func FindFilesBySourcePattern(p newpattern.SourcePattern, matcher newmatcher.MatcherInterface) (map[string]string, error) {
	m := make(map[string]string)

	if p.IsFile() {
		m[p.Path] = p.Path
		return m, nil
	}

	filepath.Walk(p.Path, func(innerPath string, info os.FileInfo, err error) error {
		if innerPath == "." || innerPath == ".." {
			return nil
		}
		//entriesWalked++
		//if reportEvery == 0 || entriesWalked % reportEvery == 0 {
		//	progressHandlerFunc(entriesWalked, entriesMatched, false)
		//}
		//
		//file := File{info, innerPath}
		//if ! filterFunc(file, err) {
		//	return nil
		//}
		//
		//entriesMatched++
		//list = append(list, file)
		//return nil

		normalizedInnerPath := filepath.ToSlash(innerPath)
		if info.IsDir() {
			normalizedInnerPath += "/"
		}

		if matcher.Matches(normalizedInnerPath) {
			m[normalizedInnerPath] = normalizedInnerPath
		}

		return nil
	})
	return m, nil
}

