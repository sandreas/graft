package filesystem

import (
	"path/filepath"
	"github.com/spf13/afero"
	"os"
	"sort"
)

func Walk(fs afero.Fs, root string, walkFn filepath.WalkFunc) error {
	info, err := lstatIfOs(fs, root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(fs, root, info, walkFn)
}

func lstatIfOs(fs afero.Fs, path string) (info os.FileInfo, err error) {
	_, ok := fs.(*afero.OsFs)
	if ok {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		info, err = os.Lstat(absPath)
	} else {
		info, err = fs.Stat(path)
	}
	return
}

func walk(fs afero.Fs, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := readDirNames(fs, path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := lstatIfOs(fs, filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

func readDirNames(fs afero.Fs, dirname string) ([]string, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func ToAbsIfOsFs(fs afero.Fs, path string) (string, error) {
	var absSrc string
	var err error

	if _, ok := fs.(*afero.OsFs); ok {
		if absSrc, err = filepath.Abs(path); err != nil {
			return path, err
		}
		return absSrc, nil
	}
	return path, nil
}
