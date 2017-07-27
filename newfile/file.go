package newfile

import (
	"github.com/sandreas/graft/newpattern"
	"os"
	"path/filepath"
)

func FindFilesBySourcePattern(p newpattern.SourcePattern) (map[string]string, error) {
	m := make(map[string]string)

	if p.IsFile() {
		m[p.Path] = p.Path
		return m, nil
	}

	compiledPattern, err := p.Compile()
	if err != nil {
		return nil, err
	}

	filepath.Walk(p.Path, func(innerPath string, info os.FileInfo, err error) error {
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

		match := compiledPattern.MatchString(normalizedInnerPath)

		if match {
			m[normalizedInnerPath] = normalizedInnerPath
		}

		return nil
	})
	return m, nil
}

