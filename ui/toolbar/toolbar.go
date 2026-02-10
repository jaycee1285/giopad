package toolbar

import (
	"path/filepath"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"giopad/app"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Toolbar is the bottom panel with vault selector and theme toggle
type Toolbar struct {
	vaultClick    widget.Clickable
	pickerClick   widget.Clickable
	themeClick    widget.Clickable
	filesClick    widget.Clickable // Mobile nav: show files
	editorClick   widget.Clickable // Mobile nav: show editor
	pathEditor    widget.Editor
	showingPath   bool
	requestFocus  bool
	vaultPath     string
	onVaultChange func(string)
	onPickFile    func() // Called when user wants to open file picker
}

// New creates a new Toolbar
func New(vaultPath string, onVaultChange func(string)) *Toolbar {
	t := &Toolbar{
		vaultPath:     vaultPath,
		onVaultChange: onVaultChange,
	}
	t.pathEditor.SingleLine = true
	t.pathEditor.Submit = true
	t.pathEditor.SetText(vaultPath)
	return t
}

// SetOnPickFile sets the callback for when user requests file picker
func (t *Toolbar) SetOnPickFile(fn func()) {
	t.onPickFile = fn
}

// SetVaultPath updates the displayed vault path
func (t *Toolbar) SetVaultPath(path string) {
	t.vaultPath = path
	t.pathEditor.SetText(path)
}

// TogglePathEditor shows/hides the vault path editor
func (t *Toolbar) TogglePathEditor() {
	t.showingPath = !t.showingPath
	if t.showingPath {
		t.pathEditor.SetText(t.vaultPath)
		t.requestFocus = true
	}
}

// Layout renders the toolbar
func (t *Toolbar) Layout(gtx C, th *material.Theme) D {
	// Handle theme click
	if t.themeClick.Clicked(gtx) {
		app.ToggleTheme()
	}

	// Handle picker click - open native file picker
	if t.pickerClick.Clicked(gtx) && t.onPickFile != nil {
		t.onPickFile()
	}

	// Handle vault click - toggle path editor
	if t.vaultClick.Clicked(gtx) {
		t.showingPath = !t.showingPath
		if t.showingPath {
			t.pathEditor.SetText(t.vaultPath)
			t.requestFocus = true
		}
	}

	// Handle path submission
	for {
		ev, ok := t.pathEditor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := ev.(widget.SubmitEvent); ok {
			newPath := t.pathEditor.Text()
			if newPath != t.vaultPath {
				t.vaultPath = newPath
				if t.onVaultChange != nil {
					t.onVaultChange(newPath)
				}
			}
			t.showingPath = false
		}
	}

	return layout.Inset{
		Top:    unit.Dp(4),
		Bottom: unit.Dp(8),
		Left:   unit.Dp(8),
		Right:  unit.Dp(8),
	}.Layout(gtx, func(gtx C) D {
		if t.showingPath {
			return t.layoutPathEditor(gtx, th)
		}
		return t.layoutButtons(gtx, th)
	})
}

func (t *Toolbar) layoutButtons(gtx C, th *material.Theme) D {
	return layout.Flex{
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		// Vault name (click to edit path manually)
		layout.Flexed(1, func(gtx C) D {
			return t.vaultClick.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top:    unit.Dp(6),
					Bottom: unit.Dp(6),
				}.Layout(gtx, func(gtx C) D {
					name := filepath.Base(t.vaultPath)
					if name == "." || name == "" {
						name = "(no vault)"
					}
					label := material.Body2(th, "📁 "+name)
					label.Color = app.Comment()
					return label.Layout(gtx)
				})
			})
		}),
		// File picker button (native dialog)
		layout.Rigid(func(gtx C) D {
			return t.pickerClick.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top:    unit.Dp(6),
					Bottom: unit.Dp(6),
					Left:   unit.Dp(8),
				}.Layout(gtx, func(gtx C) D {
					label := material.Body2(th, "[Open]")
					label.Color = app.Comment()
					return label.Layout(gtx)
				})
			})
		}),
		// Theme toggle
		layout.Rigid(func(gtx C) D {
			return t.themeClick.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top:    unit.Dp(6),
					Bottom: unit.Dp(6),
					Left:   unit.Dp(8),
				}.Layout(gtx, func(gtx C) D {
					icon := "☀"
					if !app.CurrentTheme.IsDark {
						icon = "🌙"
					}
					label := material.Body1(th, icon)
					label.Color = app.Comment()
					return label.Layout(gtx)
				})
			})
		}),
	)
}

func (t *Toolbar) layoutPathEditor(gtx C, th *material.Theme) D {
	if t.requestFocus {
		gtx.Execute(key.FocusCmd{Tag: &t.pathEditor})
		t.requestFocus = false
	}

	ed := material.Editor(th, &t.pathEditor, "Vault path...")
	ed.Color = app.Foreground()
	ed.HintColor = app.Comment()
	ed.TextSize = unit.Sp(13)
	return ed.Layout(gtx)
}

// FilesClicked returns true if the Files nav button was clicked
func (t *Toolbar) FilesClicked(gtx C) bool {
	return t.filesClick.Clicked(gtx)
}

// EditorClicked returns true if the Editor nav button was clicked
func (t *Toolbar) EditorClicked(gtx C) bool {
	return t.editorClick.Clicked(gtx)
}

// LayoutMobileNav renders the mobile bottom navigation bar
func (t *Toolbar) LayoutMobileNav(gtx C, th *material.Theme, showingEditor bool) D {
	// Handle picker click
	if t.pickerClick.Clicked(gtx) && t.onPickFile != nil {
		t.onPickFile()
	}

	// Handle theme click
	if t.themeClick.Clicked(gtx) {
		app.ToggleTheme()
	}

	return layout.Inset{
		Top:    unit.Dp(8),
		Bottom: unit.Dp(12),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Spacing:   layout.SpaceAround,
		}.Layout(gtx,
			// Files tab
			layout.Flexed(1, func(gtx C) D {
				return t.filesClick.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						label := material.Body1(th, "Files")
						if !showingEditor {
							label.Color = app.Accent()
						} else {
							label.Color = app.Comment()
						}
						return label.Layout(gtx)
					})
				})
			}),
			// Editor tab
			layout.Flexed(1, func(gtx C) D {
				return t.editorClick.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						label := material.Body1(th, "Editor")
						if showingEditor {
							label.Color = app.Accent()
						} else {
							label.Color = app.Comment()
						}
						return label.Layout(gtx)
					})
				})
			}),
			// Open vault
			layout.Flexed(1, func(gtx C) D {
				return t.pickerClick.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						label := material.Body1(th, "Open")
						label.Color = app.Comment()
						return label.Layout(gtx)
					})
				})
			}),
			// Theme toggle
			layout.Flexed(1, func(gtx C) D {
				return t.themeClick.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						icon := "Light"
						if !app.CurrentTheme.IsDark {
							icon = "Dark"
						}
						label := material.Body1(th, icon)
						label.Color = app.Comment()
						return label.Layout(gtx)
					})
				})
			}),
		)
	})
}
