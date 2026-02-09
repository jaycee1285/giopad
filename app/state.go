package app

// State holds the application state
type State struct {
	VaultPath    string   // Root path of the markdown vault
	SelectedFile string   // Currently selected file path
	OpenFiles    []string // List of open file paths (tabs)
	TreeExpanded map[string]bool // Which directories are expanded
	DrawerOpen   bool     // Mobile: is the file tree drawer open?
}

// NewState creates a new application state
func NewState(vaultPath string) *State {
	return &State{
		VaultPath:    vaultPath,
		TreeExpanded: make(map[string]bool),
	}
}
