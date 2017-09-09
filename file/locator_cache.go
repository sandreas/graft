package file

import (
	"github.com/spf13/afero"

	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sandreas/graft/filesystem"
)

type LocatorCache struct {
	Fs        afero.Fs
	cacheFile string
	Items     []string
}

func NewLocatorCache(cacheFile string) *LocatorCache {
	return &LocatorCache{
		cacheFile: cacheFile,
		Fs:        filesystem.NewOsFs(),
		Items:     []string{},
	}
}

func (i *LocatorCache) Load() error {

	_, err := i.Fs.Stat(i.cacheFile)

	if err != nil {
		return err
	}

	file, err := i.Fs.Open(i.cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	i.Items = []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i.Items = append(i.Items, strings.Trim(scanner.Text(), "\r\n"))
	}

	return nil
}

func (i *LocatorCache) Save() error {
	_, err := i.Fs.Stat(i.cacheFile)

	if !os.IsNotExist(err) {
		return errors.New("file " + i.cacheFile + " already exists")
	}

	file, err := i.Fs.Create(i.cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range i.Items {
		fmt.Fprintln(w, strings.Trim(line, "\r\n"))
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return file.Sync()
}
