package fs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Node represents a file or directory in the tree
type Node struct {
	Path     string
	Name     string
	IsDir    bool
	Depth    int
	Children []*Node
}

// MaxDepth limits recursion depth
const MaxDepth = 5

// ScanVault recursively scans a directory and returns a tree of nodes
// Filters to only .md files and directories containing .md files
func ScanVault(root string) (*Node, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	rootNode := &Node{
		Path:  root,
		Name:  info.Name(),
		IsDir: true,
		Depth: 0,
	}

	err = scanDir(rootNode, root, 0)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

func scanDir(parent *Node, path string, depth int) error {
	if depth > MaxDepth {
		return nil // Stop recursion
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Sort: directories first, then alphabetically
	sort.Slice(entries, func(i, j int) bool {
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir // directories first
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files/dirs
		if strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(path, name)

		if entry.IsDir() {
			dirNode := &Node{
				Path:  fullPath,
				Name:  name,
				IsDir: true,
				Depth: depth + 1,
			}
			// Recursively scan subdirectory
			if err := scanDir(dirNode, fullPath, depth+1); err != nil {
				continue // Skip dirs we can't read
			}
			// Only add directory if it has markdown content
			if hasMarkdownContent(dirNode) {
				parent.Children = append(parent.Children, dirNode)
			}
		} else if strings.HasSuffix(strings.ToLower(name), ".md") {
			fileNode := &Node{
				Path:  fullPath,
				Name:  name,
				IsDir: false,
				Depth: depth + 1,
			}
			parent.Children = append(parent.Children, fileNode)
		}
	}

	return nil
}

// hasMarkdownContent checks if a directory or its children contain .md files
func hasMarkdownContent(node *Node) bool {
	if !node.IsDir {
		return strings.HasSuffix(strings.ToLower(node.Name), ".md")
	}
	for _, child := range node.Children {
		if hasMarkdownContent(child) {
			return true
		}
	}
	return false
}

// FlattenTree converts a tree into a flat slice for list rendering
// Only includes expanded directories
func FlattenTree(root *Node, expanded map[string]bool) []*Node {
	var result []*Node
	flattenNode(root, expanded, &result, true)
	return result
}

func flattenNode(node *Node, expanded map[string]bool, result *[]*Node, isRoot bool) {
	if !isRoot {
		*result = append(*result, node)
	}

	if node.IsDir && (isRoot || expanded[node.Path]) {
		for _, child := range node.Children {
			flattenNode(child, expanded, result, false)
		}
	}
}
