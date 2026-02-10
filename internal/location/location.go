package location

import (
	"net/url"
	"path/filepath"
	"strings"
)

// Location represents either a local file path or a remote URL.
type Location struct {
	Path string   // Non-empty for local files (absolute path)
	URL  *url.URL // Non-nil for remote URLs
}

// FromPath creates a Location from a local file path.
func FromPath(p string) Location {
	abs, err := filepath.Abs(p)
	if err != nil {
		abs = p
	}
	return Location{Path: abs}
}

// FromURL creates a Location from a URL string.
func FromURL(u string) (Location, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return Location{}, err
	}
	return Location{URL: parsed}, nil
}

// IsLocal returns true if the location is a local file.
func (l Location) IsLocal() bool { return l.URL == nil && l.Path != "" }

// IsRemote returns true if the location is a remote URL.
func (l Location) IsRemote() bool { return l.URL != nil }

// IsZero returns true if the location is unset.
func (l Location) IsZero() bool { return l.Path == "" && l.URL == nil }

// String returns the location as a string.
func (l Location) String() string {
	if l.IsRemote() {
		return l.URL.String()
	}
	return l.Path
}

// Name returns the filename portion of the location.
func (l Location) Name() string {
	if l.IsRemote() {
		return filepath.Base(l.URL.Path)
	}
	return filepath.Base(l.Path)
}

// Dir returns the parent directory/path of the location.
func (l Location) Dir() string {
	if l.IsRemote() {
		u := *l.URL
		u.Path = filepath.Dir(u.Path)
		return u.String()
	}
	return filepath.Dir(l.Path)
}

// IsLikelyURL checks if a string looks like a URL.
func IsLikelyURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// IsMaybeMarkdown checks if a filename has a markdown extension.
func IsMaybeMarkdown(s string) bool {
	ext := strings.ToLower(filepath.Ext(s))
	return ext == ".md" || ext == ".markdown"
}
