package sftpd

import (
	"path/filepath"
	"sort"
	"strings"
	"errors"
	"os"
)

type PathMapper struct {
	tree map[string][]string
	basePath string
}

func NewPathMapper(files []string, basePath string) *PathMapper {
	pathMapper := &PathMapper{
		basePath: basePath,
	}
	pathMapper.normalizeBasePath()
	pathMapper.buildTree(files)
	return pathMapper
}

func (mapper *PathMapper) normalizeBasePath() {
	mapper.basePath = mapper.normalizePath(mapper.basePath)
	if mapper.basePath == "" {
		mapper.basePath = "."
	}
}

func (mapper *PathMapper) normalizePath(basePath string) string {
	basePath = filepath.ToSlash(basePath)
	basePath = strings.TrimPrefix(basePath, "./")
	if basePath == "." {
		return ""
	}
	return strings.TrimRight(basePath, "/")
}

func (mapper *PathMapper) List(key string) ([]string, bool) {
	normalizedKey := mapper.slashify(key)
	value, ok := mapper.tree[normalizedKey]
	return value, ok
}

func (mapper *PathMapper) PathTo(reference string) (string, error) {
	normalizedKey := mapper.slashify(reference)
	_, ok := mapper.tree[normalizedKey]
	if ! ok {
		return "", errors.New("PathTo " + reference + " not found")
	}
	return filepath.FromSlash(mapper.normalizePath(mapper.basePath) + "/" + reference), nil
}

func (mapper *PathMapper) Stat(reference string) (os.FileInfo, error) {
	path, err := mapper.PathTo(reference)
	if err != nil {
		return nil, err
	}
	return os.Stat(path)
}

func (mapper *PathMapper) slashify(path string) string {
	toSlash := filepath.ToSlash(path)
	trimmed := strings.TrimLeft(toSlash, "/")
	return "/" + trimmed
}


func (mapper *PathMapper) buildTree(matchingPaths []string) {
	mapper.tree = make(map[string][]string)

	sort.Strings(matchingPaths)

	for _, path := range matchingPaths {
		normalizedPath := mapper.normalizePath(path)
		key, parentPath := mapper.normalizePathMapItem(normalizedPath)

		for {
			pathToAppend := strings.TrimPrefix(mapper.normalizePath(path), mapper.basePath)
			mapper.tree[key] = append(mapper.tree[key], pathToAppend)
			path = parentPath
			key, parentPath = mapper.normalizePathMapItem(parentPath)
			if _, ok := mapper.tree[key]; ok {
				break
			}
		}
	}
}

func (mapper *PathMapper) normalizePathMapItem(path string) (string, string) {
	parentPath := mapper.normalizePath(filepath.Dir(path))
	parentWithoutBaseDir := strings.TrimPrefix(parentPath, mapper.basePath)
	key := mapper.slashify(parentWithoutBaseDir)
	return key, parentPath
}

