package sftpd

import (
	"path/filepath"
	"strings"
)

type FileTree struct {
	basePath string
	key      string
	children []FileTree
}


func normalizePath(path string) string {
	path = filepath.ToSlash(path)
	path = strings.TrimPrefix(path, "./")
	if path == "." {
		return ""
	}
	return strings.TrimRight(path, "/")
}


func NewFileTree(basePath string) *FileTree {
	basePath = normalizePath(basePath)
	rootPath := ""
	if strings.HasPrefix(basePath, "/") {
		rootPath = "/"
	}
	return &FileTree{
		basePath: basePath,
		key: rootPath,
	}
}

//func (tree *FileTree) Integrate(path string) error {
//	path = strings.TrimPrefix(normalizePath(path), tree.basePath)
//
//	if tree.key == "" && strings.HasPrefix(path, "/") {
//		return errors.New("Absolute path " + path + " cannot be integrated into relative tree")
//	}
//
//	parts := strings.Split(path, "/")
//
//	len := len(tree.children)
//	for i:=0;i<len;i++ {
//		if err := tree.children[i].SubIntegrate(path); err == nil {
//			return nil
//		}
//	}
//
//	return nil
//}


func (tree *FileTree) SubIntegrate(path string) bool {
	comparePath := "/" + strings.TrimPrefix("/", path)
	return comparePath == tree.key || strings.HasPrefix(comparePath, tree.key + "/")
}
