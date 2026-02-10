// SPDX-License-Identifier: Unlicense OR MIT

//go:build !android

package fs

// SAFEntry represents a file or directory in SAF
type SAFEntry struct {
	Name  string
	URI   string
	IsDir bool
}

// ListSAFDir is not available on non-Android platforms
func ListSAFDir(treeURI string) ([]SAFEntry, error) {
	return nil, nil
}

// ListSAFSubDir is not available on non-Android platforms
func ListSAFSubDir(treeURI, docURI string) ([]SAFEntry, error) {
	return nil, nil
}

// ReadSAFFile is not available on non-Android platforms
func ReadSAFFile(docURI string) ([]byte, error) {
	return nil, nil
}

// WriteSAFFile is not available on non-Android platforms
func WriteSAFFile(docURI string, data []byte) error {
	return nil
}

// GetSAFTreeName is not available on non-Android platforms
func GetSAFTreeName(treeURI string) string {
	return ""
}

// IsSAFURI checks if a path is a SAF content URI
func IsSAFURI(path string) bool {
	return false
}
