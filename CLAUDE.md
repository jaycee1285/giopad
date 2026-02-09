# Claude Instructions for giopad

## Project Overview

Go + Gio markdown editor with file tree. Targets Linux desktop and Android. Learning project for Gio framework - building something between crustdown (TUI) and Ferrite (full egui editor).

**Stack:** Go, Gio (gioui.org), gioui.org/x/markdown
**Platforms:** Linux (Wayland/X11), Android
**Packaging:** Nix flake, gogio for Android

## Key Paths

| Path | Purpose |
|------|---------|
| `main.go` | Entry point, window setup |
| `app/theme.go` | Ayu-Mirage color palette |
| `app/state.go` | Application state |
| `ui/tree/` | File tree widget |
| `ui/editor/` | Markdown editor/viewer |
| `ui/layout/` | Desktop vs mobile layout |
| `fs/` | Vault operations, file watching |

## Build Commands

```bash
# Desktop
go build -o giopad .

# Android
gogio -target android -appid io.github.jaycee1285.giopad .

# With Nix
nix build
nix develop  # Enter dev shell
```

## Design Decisions

- **Immediate mode GUI**: Gio is immediate mode like egui/imgui. Redraw every frame.
- **No webview**: Pure Go, no JS, no CSS. Theming via Go color structs.
- **x/markdown**: Use gioui.org/x/markdown for rendering, not custom parser.
- **goldmark**: x/markdown uses goldmark internally for parsing.

## Conventions

- Use `layout.Context` alias `C` and `layout.Dimensions` alias `D` for brevity
- Theme colors in `app/theme.go`
- Platform-specific layouts in `ui/layout/`

## Owner

- GitHub: jaycee1285
- User: john
