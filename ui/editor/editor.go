package editor

import (
	"os"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/markdown"
	"gioui.org/x/richtext"

	"giopad/app"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Editor displays markdown content
type Editor struct {
	currentPath string
	content     []byte
	renderer    *markdown.Renderer
	spans       []richtext.SpanStyle
	textState   richtext.InteractiveText
	list        layout.List

	// Edit mode
	editMode   bool
	textEditor widget.Editor
	dirty      bool
}

// New creates a new Editor
func New() *Editor {
	r := markdown.NewRenderer()
	// Configure colors
	r.Config.DefaultColor = app.Foreground()
	r.Config.InteractiveColor = app.Blue()
	r.Config.DefaultSize = unit.Sp(14)
	r.Config.H1Size = unit.Sp(28)
	r.Config.H2Size = unit.Sp(24)
	r.Config.H3Size = unit.Sp(20)
	r.Config.H4Size = unit.Sp(16)

	e := &Editor{
		renderer: r,
		list:     layout.List{Axis: layout.Vertical},
	}
	e.textEditor.SingleLine = false
	e.textEditor.Submit = false
	return e
}

// LoadFile loads and parses a markdown file
func (e *Editor) LoadFile(path string) error {
	if path == e.currentPath {
		return nil // Already loaded
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	e.currentPath = path
	e.content = content
	e.editMode = false
	e.dirty = false

	// Set editor content
	e.textEditor.SetText(string(content))

	// Parse markdown to richtext spans
	spans, err := e.renderer.Render(content)
	if err != nil {
		return err
	}
	e.spans = spans

	return nil
}

// ToggleEdit switches between view and edit mode
func (e *Editor) ToggleEdit() {
	if e.currentPath == "" {
		return
	}

	if e.editMode {
		// Leaving edit mode - update content and re-render
		newContent := e.textEditor.Text()
		if newContent != string(e.content) {
			e.content = []byte(newContent)
			e.dirty = true
			// Re-render markdown
			if spans, err := e.renderer.Render(e.content); err == nil {
				e.spans = spans
			}
		}
	} else {
		// Entering edit mode - focus at start
		e.textEditor.SetCaret(0, 0)
	}
	e.editMode = !e.editMode
}

// Save writes content to disk
func (e *Editor) Save() error {
	if e.currentPath == "" || !e.dirty {
		return nil
	}
	err := os.WriteFile(e.currentPath, e.content, 0644)
	if err == nil {
		e.dirty = false
	}
	return err
}

// IsDirty returns true if there are unsaved changes
func (e *Editor) IsDirty() bool {
	return e.dirty
}

// IsEditMode returns true if in edit mode
func (e *Editor) IsEditMode() bool {
	return e.editMode
}

// Layout renders the markdown content
func (e *Editor) Layout(gtx C, th *material.Theme) D {
	if e.currentPath == "" {
		label := material.Body1(th, "Select a file")
		label.Color = app.Comment()
		return layout.Center.Layout(gtx, label.Layout)
	}

	return layout.Inset{
		Top:    unit.Dp(16),
		Left:   unit.Dp(24),
		Right:  unit.Dp(24),
		Bottom: unit.Dp(16),
	}.Layout(gtx, func(gtx C) D {
		if e.editMode {
			// Edit mode - raw text editor
			ed := material.Editor(th, &e.textEditor, "")
			ed.Color = app.Foreground()
			ed.HintColor = app.Comment()
			ed.TextSize = unit.Sp(14)
			ed.LineHeight = unit.Sp(20)
			ed.Editor.Alignment = text.Start
			return e.list.Layout(gtx, 1, func(gtx C, _ int) D {
				return ed.Layout(gtx)
			})
		}

		// View mode - rendered markdown
		if len(e.spans) == 0 {
			label := material.Body1(th, "(empty file)")
			label.Color = app.Comment()
			return layout.Center.Layout(gtx, label.Layout)
		}

		return e.list.Layout(gtx, 1, func(gtx C, _ int) D {
			return richtext.Text(&e.textState, th.Shaper, e.spans...).Layout(gtx)
		})
	})
}

// CurrentPath returns the currently loaded file path
func (e *Editor) CurrentPath() string {
	return e.currentPath
}
