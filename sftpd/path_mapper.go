package sftpd

import (
	"path/filepath"
	"sort"
	"strings"
)

type PathMapper struct {
	tree map[string][]string
	basePath string
}

func NewPathMapper(files []string, basePath string) *PathMapper {
	pathMapper := &PathMapper{
		basePath: basePath,
	}
	pathMapper.buildTree(files)
	return pathMapper
}

func (mapper *PathMapper) Get(key string) ([]string, bool) {

	normalizedKey := mapper.slashify(key)

	value, ok := mapper.tree[normalizedKey]
	return value, ok
}

func (mapper *PathMapper) slashify(path string) string {
	toSlash := filepath.ToSlash(path)
	trimmed := strings.TrimLeft(toSlash, "/")
	return "/" + trimmed
}


func (mapper *PathMapper) buildTree(matchingPaths []string) {
	mapper.tree = make(map[string][]string)

	sort.Strings(matchingPaths)

	//if val, ok := dict["foo"]; ok {
	//	//do something here
	//}

	for _, path := range matchingPaths {
		key, parentPath := mapper.normalizePathMapItem(path)

		for {
			// println("append: ", key, " => ", path)
			mapper.tree[key] = append(mapper.tree[key], path)
			path = parentPath
			//println("before => key:", key, "parentPath:", parentPath)
			key, parentPath = mapper.normalizePathMapItem(parentPath)
			//println("after  => key:", key, "parentPath:", parentPath)
			_, ok := mapper.tree[key]

			//println("is present?", key, ok)
			if ok {
				break
			}
		}
	}
}

func (mapper *PathMapper) normalizePathMapItem(path string) (string, string) {
	parentPath := filepath.ToSlash(filepath.Dir(path))
	key := parentPath
	if parentPath == "." {
		key = "/"
	}

	key = mapper.slashify(key)

	return key, parentPath
}
