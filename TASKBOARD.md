# giopad Taskboard

## Current Sprint

### In Progress
- [ ] **Android black screen bug** - App builds but shows black screen on device

### Up Next
- [ ] Debug Android rendering with logcat

---

## Known Issues

- **Android black screen**: APK installs and runs but shows only black screen. Attempted fixes:
  - Removed hardcoded vault path that doesn't exist on Android
  - Made explorer initialization lazy
  - Added fallback for zero constraints
  - Forced mobile layout detection for narrow screens
  - None resolved the issue

### Debug Steps Needed
1. **Human**: Run `adb logcat | grep -i gio` while launching app
2. **Human**: Check for crash logs or error messages
3. **Computer**: Analyze logs and fix root cause

## Vendor Patches

- **gioui.org/x/explorer**: Added `ChooseDirectory()` for Linux (xdg-portal) and Android (ACTION_OPEN_DOCUMENT_TREE). Ready for upstream.
- **gioui.org/x/explorer**: Added SAF file operations (listDir, listSubDir, readFile, writeFile, getTreeName) for Android.
- **gioui.org/x/markdown**: Added soft/hard line break handling in `renderText` (3 lines). Could upstream.

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
- [ ] SAF file traversal via JNI (DocumentFile bridge)
- [ ] Back button handling
- [ ] Touch-friendly sizing

---

## Completed (2026-02-09)
- [x] Dirty indicator in window title (real-time)
- [x] Toolbar with vault selector + theme toggle
- [x] Fix line break rendering (vendor patch)
- [x] Keyboard shortcuts: Ctrl+T (theme), Ctrl+D (directory), Ctrl+O (open file), Ctrl+Shift+O (vault picker)
- [x] Fix editor auto-focus on Ctrl+E and Ctrl+D
- [x] Native file picker via x/explorer (Linux xdg-portal, Android SAF)
- [x] Native directory picker via x/explorer (Linux + Android)
- [x] Toolbar [Open] button for vault selection
- [x] Shared location package from crustdown (internal/location)
- [x] Android build via gogio (stack-build script)
- [x] Android directory picker (ACTION_OPEN_DOCUMENT_TREE)
- [x] SAF file traversal (fs/saf_android.go - JNI bridge to DocumentFile)
- [x] SAF file read/write support
- [x] Mobile layout code (tree OR editor with bottom nav)
- [x] fs.ReadFile/WriteFile with SAF support

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

*Last updated: 2026-02-09*
