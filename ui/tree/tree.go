package tree

import (
	"image"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"giopad/app"
	"giopad/fs"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Tree is the file tree widget
type Tree struct {
	Root     *fs.Node
	Expanded map[string]bool
	Selected string
	Focused  bool

	list      widget.List
	clicks    map[string]*widget.Clickable
	flatNodes []*fs.Node // cached for keyboard nav
}

// New creates a new Tree widget
func New() *Tree {
	return &Tree{
		Expanded: make(map[string]bool),
		clicks:   make(map[string]*widget.Clickable),
	}
}

// SetRoot sets the root node and rescans
func (t *Tree) SetRoot(root *fs.Node) {
	t.Root = root
	// Expand root by default
	if root != nil {
		t.Expanded[root.Path] = true
		// Select first item if nothing selected
		if t.Selected == "" && len(root.Children) > 0 {
			t.Selected = root.Children[0].Path
		}
	}
}

// clickable returns or creates a clickable for a path
func (t *Tree) clickable(path string) *widget.Clickable {
	if c, ok := t.clicks[path]; ok {
		return c
	}
	c := new(widget.Clickable)
	t.clicks[path] = c
	return c
}

// Layout renders the tree
func (t *Tree) Layout(gtx C, th *material.Theme) D {
	if t.Root == nil {
		return material.Body1(th, "No vault loaded").Layout(gtx)
	}

	nodes := fs.FlattenTree(t.Root, t.Expanded)
	t.flatNodes = nodes // cache for keyboard nav

	// Handle keyboard navigation when focused
	if t.Focused {
		t.handleKeys(gtx)
	}

	t.list.Axis = layout.Vertical

	return material.List(th, &t.list).Layout(gtx, len(nodes), func(gtx C, i int) D {
		node := nodes[i]
		return t.layoutNode(gtx, th, node)
	})
}

func (t *Tree) handleKeys(gtx C) {
	for {
		ev, ok := gtx.Event(
			key.Filter{Name: key.NameDownArrow},
			key.Filter{Name: key.NameUpArrow},
			key.Filter{Name: key.NameReturn},
			key.Filter{Name: key.NameSpace},
			key.Filter{Name: key.NameRightArrow},
			key.Filter{Name: key.NameLeftArrow},
		)
		if !ok {
			break
		}
		e, ok := ev.(key.Event)
		if !ok || e.State != key.Press {
			continue
		}

		idx := t.selectedIndex()

		switch e.Name {
		case key.NameDownArrow:
			if idx < len(t.flatNodes)-1 {
				t.Selected = t.flatNodes[idx+1].Path
			}
		case key.NameUpArrow:
			if idx > 0 {
				t.Selected = t.flatNodes[idx-1].Path
			}
		case key.NameReturn, key.NameSpace:
			// Toggle dir or select file
			if idx >= 0 && idx < len(t.flatNodes) {
				node := t.flatNodes[idx]
				if node.IsDir {
					t.Expanded[node.Path] = !t.Expanded[node.Path]
				}
			}
		case key.NameRightArrow:
			// Expand directory
			if idx >= 0 && idx < len(t.flatNodes) {
				node := t.flatNodes[idx]
				if node.IsDir {
					t.Expanded[node.Path] = true
				}
			}
		case key.NameLeftArrow:
			// Collapse directory
			if idx >= 0 && idx < len(t.flatNodes) {
				node := t.flatNodes[idx]
				if node.IsDir && t.Expanded[node.Path] {
					t.Expanded[node.Path] = false
				}
			}
		}
	}
}

func (t *Tree) selectedIndex() int {
	for i, n := range t.flatNodes {
		if n.Path == t.Selected {
			return i
		}
	}
	return -1
}

func (t *Tree) layoutNode(gtx C, th *material.Theme, node *fs.Node) D {
	click := t.clickable(node.Path)

	// Handle clicks
	if click.Clicked(gtx) {
		if node.IsDir {
			// Toggle expansion
			t.Expanded[node.Path] = !t.Expanded[node.Path]
		} else {
			// Select file
			t.Selected = node.Path
		}
	}

	// Indent based on depth
	indent := unit.Dp(float32(node.Depth) * 16)

	// Selection highlight
	isSelected := t.Selected == node.Path

	return click.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			// Background (for selection)
			layout.Expanded(func(gtx C) D {
				if isSelected {
					rect := image.Rectangle{Max: gtx.Constraints.Min}
					paint.FillShape(gtx.Ops, app.Selection(), clip.Rect(rect).Op())
				}
				return D{Size: gtx.Constraints.Min}
			}),
			// Content
			layout.Stacked(func(gtx C) D {
				return layout.Inset{
					Left:   indent,
					Top:    unit.Dp(4),
					Bottom: unit.Dp(4),
					Right:  unit.Dp(8),
				}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						// Icon
						layout.Rigid(func(gtx C) D {
							icon := t.nodeIcon(node)
							label := material.Body2(th, icon)
							label.Color = app.Comment()
							return layout.Inset{Right: unit.Dp(6)}.Layout(gtx, label.Layout)
						}),
						// Name
						layout.Flexed(1, func(gtx C) D {
							label := material.Body2(th, node.Name)
							if node.IsDir {
								label.Color = app.Blue()
							} else {
								label.Color = app.Foreground()
							}
							return label.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

func (t *Tree) nodeIcon(node *fs.Node) string {
	if node.IsDir {
		if t.Expanded[node.Path] {
			return "▼"
		}
		return "▶"
	}
	return "•"
}

// SelectedPath returns the currently selected file path
func (t *Tree) SelectedPath() string {
	return t.Selected
}
