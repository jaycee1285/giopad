package fs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ReadFile reads a file from either filesystem or SAF URI
func ReadFile(path string) ([]byte, error) {
	if IsSAFURI(path) {
		return ReadSAFFile(path)
	}
	return os.ReadFile(path)
}

// WriteFile writes to a file, using SAF on Android for content:// URIs
func WriteFile(path string, data []byte) error {
	if IsSAFURI(path) {
		return WriteSAFFile(path, data)
	}
	return os.WriteFile(path, data, 0644)
}

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
// Supports both filesystem paths and Android SAF content:// URIs
func ScanVault(root string) (*Node, error) {
	// Handle Android SAF URIs
	if IsSAFURI(root) {
		return scanVaultSAF(root)
	}

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

// scanVaultSAF scans a SAF tree URI on Android
func scanVaultSAF(treeURI string) (*Node, error) {
	name := GetSAFTreeName(treeURI)
	if name == "" {
		name = "Vault"
	}

	rootNode := &Node{
		Path:  treeURI,
		Name:  name,
		IsDir: true,
		Depth: 0,
	}

	err := scanSAFDir(rootNode, treeURI, treeURI, 0)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

func scanSAFDir(parent *Node, treeURI, docURI string, depth int) error {
	if depth > MaxDepth {
		return nil
	}

	var entries []SAFEntry
	var err error

	if depth == 0 {
		entries, err = ListSAFDir(treeURI)
	} else {
		entries, err = ListSAFSubDir(treeURI, docURI)
	}
	if err != nil {
		return err
	}

	// Sort: directories first, then alphabetically
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name, ".") {
			continue
		}

		if entry.IsDir {
			dirNode := &Node{
				Path:  entry.URI,
				Name:  entry.Name,
				IsDir: true,
				Depth: depth + 1,
			}
			if err := scanSAFDir(dirNode, treeURI, entry.URI, depth+1); err != nil {
				continue
			}
			if hasMarkdownContent(dirNode) {
				parent.Children = append(parent.Children, dirNode)
			}
		} else if strings.HasSuffix(strings.ToLower(entry.Name), ".md") {
			fileNode := &Node{
				Path:  entry.URI,
				Name:  entry.Name,
				IsDir: false,
				Depth: depth + 1,
			}
			parent.Children = append(parent.Children, fileNode)
		}
	}

	return nil
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
