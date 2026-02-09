# giopad Taskboard

## Current Sprint

### In Progress
- [ ] Fix editor auto-focus on Ctrl+E (need to find correct Gio focus API)

### Up Next
- [ ] Dirty indicator (show * in title or UI when unsaved)
- [ ] File picker / vault selector
- [ ] Mobile layout (drawer instead of side panel)
- [ ] Android build testing with gogio

---

## Known Issues

- **Editor focus**: Pressing Ctrl+E enters edit mode but cursor doesn't appear until you click. The Gio focus API changed and we need to find the correct incantation.

---

## Backlog

### Core Features
- [ ] Search within vault
- [ ] Recent files
- [ ] File creation (new file)
- [ ] Delete file (with confirmation)

### Polish
- [ ] File watcher for external changes
- [ ] More keyboard shortcuts (Ctrl+W close, etc.)
- [ ] Tab support for multiple open files
- [ ] Syntax highlighting in edit mode
- [ ] Scroll position preservation when switching files

### Android-Specific
- [ ] SAF (Storage Access Framework) for vault selection
- [ ] Back button handling
- [ ] Touch-friendly sizing
- [ ] Swipe to open/close tree drawer

---

## Completed (2026-02-08)
- [x] Project scaffold with Nix flake
- [x] Ayu-Mirage theme
- [x] File tree widget with expand/collapse
- [x] Markdown rendering via x/markdown
- [x] Scrollable content
- [x] Ayu Light theme + live toggle (sun/moon button)
- [x] Edit mode with raw text editing (Ctrl+E)
- [x] Save functionality (Ctrl+S)
- [x] Keyboard navigation in tree (up/down/left/right/enter)
- [x] Focus switching (Ctrl+Left to tree, Ctrl+Right to editor)
- [x] Initial selection in tree on launch

---

## Stats

- **Lines of Go**: ~750 (with comments/whitespace)
- **Build time**: ~5-10 seconds
- **Dependencies**: gioui.org, gioui.org/x

---

*Last updated: 2026-02-08*
