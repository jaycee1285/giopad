package main

import (
	"image"
	"io"
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
	"gioui.org/x/explorer"

	appstate "giopad/app"
	"giopad/fs"
	"giopad/ui/editor"
	"giopad/ui/toolbar"
	"giopad/ui/tree"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	log.Println("giopad: main() starting")
	go func() {
		log.Println("giopad: creating window")
		w := new(app.Window)
		w.Option(app.Title("giopad"))
		w.Option(app.Size(1200, 800))
		log.Println("giopad: window created, calling run()")

		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	log.Println("giopad: calling app.Main()")
	app.Main()
}

// FileOpenResult carries the result of a file open operation.
type FileOpenResult struct {
	Path    string
	Content string
	Err     error
}

// VaultPickResult carries the result of a vault picker operation.
type VaultPickResult struct {
	VaultPath string
	Err       error
}

func run(w *app.Window) error {
	log.Println("giopad: run() started")
	var ops op.Ops
	var lastTitle string

	// Mobile layout state: show tree (0) or editor (1)
	showingEditor := false
	const mobileBreakpoint = 600 // dp

	// Initialize tree and editor
	log.Println("giopad: initializing tree and editor")
	fileTree := tree.New()
	mdEditor := editor.New()
	log.Println("giopad: tree and editor initialized")

	// File explorer - initialized lazily
	var expl *explorer.Explorer
	fileOpenCh := make(chan FileOpenResult, 1)
	vaultPickCh := make(chan VaultPickResult, 1)

	getExplorer := func() *explorer.Explorer {
		if expl == nil {
			expl = explorer.NewExplorer(w)
		}
		return expl
	}

	// Default vault path - empty until user picks one on mobile
	vaultPath := ""

	// Scan vault
	scanVault := func(path string) {
		if path == "" {
			return
		}
		if root, err := fs.ScanVault(path); err == nil {
			fileTree.SetRoot(root)
		}
	}

	// On desktop, try default path
	if home, err := os.UserHomeDir(); err == nil {
		defaultPath := filepath.Join(home, "Sync", "JMC", "SideProjects")
		if _, err := os.Stat(defaultPath); err == nil {
			vaultPath = defaultPath
			scanVault(vaultPath)
		}
	}

	// Initialize toolbar with vault change callback
	bottomBar := toolbar.New(vaultPath, func(newPath string) {
		vaultPath = newPath
		scanVault(newPath)
	})

	// pickVault launches native directory picker
	pickVault := func() {
		go func() {
			dirPath, err := getExplorer().ChooseDirectory()
			if err != nil {
				if err != explorer.ErrUserDecline {
					vaultPickCh <- VaultPickResult{Err: err}
				}
				return
			}
			if dirPath != "" {
				vaultPickCh <- VaultPickResult{VaultPath: dirPath}
				w.Invalidate()
			}
		}()
	}

	// Wire up toolbar's file picker button
	bottomBar.SetOnPickFile(pickVault)
	log.Println("giopad: initialization complete, entering event loop")

	// openFile launches the native file picker in a goroutine
	openFile := func() {
		go func() {
			file, err := getExplorer().ChooseFile(".md", ".markdown")
			if err != nil {
				if err != explorer.ErrUserDecline {
					fileOpenCh <- FileOpenResult{Err: err}
				}
				return
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				fileOpenCh <- FileOpenResult{Err: err}
				return
			}

			// Try to get the file path if available
			var path string
			if f, ok := file.(*os.File); ok {
				path = f.Name()
			}

			fileOpenCh <- FileOpenResult{Path: path, Content: string(content)}
			w.Invalidate()
		}()
	}

	for {
		// Check for file open results (non-blocking)
		select {
		case result := <-fileOpenCh:
			if result.Err != nil {
				log.Printf("file open error: %v", result.Err)
			} else if result.Path != "" {
				mdEditor.LoadFile(result.Path)
			}
		case result := <-vaultPickCh:
			if result.Err != nil {
				log.Printf("vault pick error: %v", result.Err)
			} else if result.VaultPath != "" {
				vaultPath = result.VaultPath
				scanVault(vaultPath)
				bottomBar.SetVaultPath(vaultPath)
			}
		default:
		}

		ev := w.Event()
		if expl != nil {
			expl.ListenEvents(ev)
		}

		switch e := ev.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// Log first few frames to confirm rendering
			if lastTitle == "" {
				log.Printf("giopad: first FrameEvent, constraints=%v, metric=%+v", gtx.Constraints, gtx.Metric)
			}
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
			for {
				ev, ok := gtx.Event(key.Filter{Name: "T", Required: key.ModCtrl})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					appstate.ToggleTheme()
				}
			}
			for {
				ev, ok := gtx.Event(key.Filter{Name: "D", Required: key.ModCtrl})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					bottomBar.TogglePathEditor()
				}
			}
			for {
				ev, ok := gtx.Event(key.Filter{Name: "O", Required: key.ModCtrl})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					openFile()
				}
			}
			// Ctrl+Shift+O for vault picker
			for {
				ev, ok := gtx.Event(key.Filter{Name: "O", Required: key.ModCtrl | key.ModShift})
				if !ok {
					break
				}
				if e, ok := ev.(key.Event); ok && e.State == key.Press {
					pickVault()
				}
			}

			// Update window title with dirty indicator
			title := "giopad"
			if path := mdEditor.CurrentPath(); path != "" {
				name := filepath.Base(path)
				if mdEditor.IsDirty() {
					name += "*"
				}
				title = name + " - giopad"
			}
			if title != lastTitle {
				w.Option(app.Title(title))
				lastTitle = title
			}

			// Fill background (use white if constraints seem wrong)
			bg := appstate.Background()
			maxPt := gtx.Constraints.Max
			if maxPt.X <= 0 || maxPt.Y <= 0 {
				maxPt = image.Pt(1000, 2000) // fallback
			}
			paint.FillShape(gtx.Ops, bg, clip.Rect(image.Rectangle{Max: maxPt}).Op())

			// Detect mobile layout (width < 600dp)
			// Force mobile on Android or narrow screens
			screenWidthDp := float32(gtx.Constraints.Max.X) / gtx.Metric.PxPerDp
			isMobile := screenWidthDp < float32(mobileBreakpoint) || gtx.Constraints.Max.X < 800
			// Log layout decision once
			if lastTitle == "" {
				log.Printf("giopad: screenWidthDp=%.1f, isMobile=%v, maxX=%d", screenWidthDp, isMobile, gtx.Constraints.Max.X)
			}

			// Handle file selection - switch to editor on mobile
			selected := fileTree.SelectedPath()
			if selected != "" && selected != mdEditor.CurrentPath() {
				mdEditor.LoadFile(selected)
				if isMobile {
					showingEditor = true
				}
			}

			if isMobile {
				// Mobile: show tree OR editor, with bottom nav
				// Handle nav button clicks before layout
				if bottomBar.FilesClicked(gtx) {
					showingEditor = false
				}
				if bottomBar.EditorClicked(gtx) {
					showingEditor = true
				}

				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Main content area
					layout.Flexed(1, func(gtx C) D {
						// Surface background for content
						paint.FillShape(gtx.Ops, appstate.Surface(), clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())

						if showingEditor {
							return mdEditor.Layout(gtx, th)
						}
						// Tree view with padding
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
							return fileTree.Layout(gtx, th)
						})
					}),
					// Bottom nav bar
					layout.Rigid(func(gtx C) D {
						return bottomBar.LayoutMobileNav(gtx, th, showingEditor)
					}),
				)
			} else {
				// Desktop: side-by-side layout
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
							// Toolbar at bottom
							layout.Rigid(func(gtx C) D {
								return bottomBar.Layout(gtx, th)
							}),
						)
					}),
					// Content area
					layout.Flexed(1, func(gtx C) D {
						return mdEditor.Layout(gtx, th)
					}),
				)
			}

			e.Frame(gtx.Ops)
		}
	}
}
