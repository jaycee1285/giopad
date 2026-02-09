package main

import (
	"image"
	"log"
	"os"
	"path/filepath"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	appstate "giopad/app"
	"giopad/fs"
	"giopad/ui/editor"
	"giopad/ui/toggle"
	"giopad/ui/tree"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("giopad"))
		w.Option(app.Size(1200, 800))

		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops

	// Initialize tree, editor, and theme toggle
	fileTree := tree.New()
	mdEditor := editor.New()
	themeToggle := toggle.New()

	// Default vault path - can be changed later
	home, _ := os.UserHomeDir()
	vaultPath := filepath.Join(home, "Sync", "JMC", "SideProjects")

	// Scan vault
	if root, err := fs.ScanVault(vaultPath); err == nil {
		fileTree.SetRoot(root)
	}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			th := appstate.AyuMirageTheme() // Recreate each frame to pick up theme changes

			// Handle key events
			for {
				ev, ok := gtx.Event(key.Filter{Name: "E", Required: key.ModCtrl})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					mdEditor.ToggleEdit()
				}
			}
			for {
				ev, ok := gtx.Event(key.Filter{Name: "S", Required: key.ModCtrl})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					mdEditor.Save()
				}
			}
			for {
				ev, ok := gtx.Event(key.Filter{Name: key.NameLeftArrow, Required: key.ModCtrl})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					fileTree.Focused = true
				}
			}
			for {
				ev, ok := gtx.Event(key.Filter{Name: key.NameRightArrow, Required: key.ModCtrl})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					fileTree.Focused = false
				}
			}

			// Fill background
			paint.FillShape(gtx.Ops, appstate.Background(), clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())

			// Layout: tree on left, content placeholder on right
			layout.Flex{}.Layout(gtx,
				// Tree panel
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Max.X = 280
					gtx.Constraints.Min.X = 280

					// Tree panel background
					paint.FillShape(gtx.Ops, appstate.Surface(), clip.Rect(image.Rectangle{Max: image.Pt(280, gtx.Constraints.Max.Y)}).Op())

					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						// Tree takes most space
						layout.Flexed(1, func(gtx C) D {
							return layout.Inset{
								Top:    unit.Dp(8),
								Left:   unit.Dp(8),
								Right:  unit.Dp(8),
								Bottom: unit.Dp(8),
							}.Layout(gtx, func(gtx C) D {
								return fileTree.Layout(gtx, th)
							})
						}),
						// Theme toggle at bottom
						layout.Rigid(func(gtx C) D {
							return themeToggle.Layout(gtx, th)
						}),
					)
				}),
				// Content area
				layout.Flexed(1, func(gtx C) D {
					// Load selected file into editor
					selected := fileTree.SelectedPath()
					if selected != "" && selected != mdEditor.CurrentPath() {
						mdEditor.LoadFile(selected)
					}
					return mdEditor.Layout(gtx, th)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
