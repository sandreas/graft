package file

import (
	"path/filepath"
	"strings"
	"os"
)


type TreeNode struct {
	basePath string
	key      string
	children map[string]*TreeNode
}


func NewTreeNode(basePath string) *TreeNode {
	return &TreeNode{
		basePath:       basePath,
		children: make(map[string]*TreeNode, 0),
	}
}

func (t *TreeNode) Integrate(path string) error {
	cleanedPath := filepath.Clean(path)
	splitted := strings.Split(cleanedPath, string(os.PathSeparator))
	t.integrateSplitted(splitted)
	return nil
}

func (t *TreeNode) integrateSplitted(splittedPath []string) error {
	if len(splittedPath) == 0 {
		return nil
	}
	key := splittedPath[0]
	node, ok := t.children[key]
	if !ok {
		node = &TreeNode{
			key: key,
			children: make(map[string]*TreeNode, 0),
		}
	}
	t.children[key] = node
	return node.integrateSplitted(splittedPath[1:])
}

func (t *TreeNode) List(path string) []string {
	cleanedPath := filepath.Clean(path)
	splitted := strings.Split(cleanedPath, string(os.PathSeparator))

	return t.listSplitted(splitted)
}

func (t *TreeNode) listSplitted(splittedPath []string) []string {
	if len(splittedPath) == 0 {
		keys := make([]string, 0, len(t.children))
		for k := range t.children {
			keys = append(keys, k)
		}
		return keys
	}

	key := splittedPath[0]
	node, ok := t.children[key]
	if !ok {
		return []string{}
	}
	return node.listSplitted(splittedPath[1:])
}