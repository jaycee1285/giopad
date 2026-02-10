# Shared Core: crustdown + giopad

## The Insight

From product review: "The event handling in BubbleTea apps is basically the same as GUI apps - it's just the rendering that differs. Can we reuse crustdown's logic in giopad?"

**Verdict: Yes, and it's worth doing.**

This isn't just architecturally elegant - it's practical. crustdown has ~6 months of edge case fixes baked into its tree navigation, clipboard handling, and file operations. Reimplementing that in giopad means re-discovering those edge cases.

## What's Actually Portable

### Tier 1: Copy tomorrow (zero changes needed)

| Package | Lines | What it does |
|---------|-------|--------------|
| `location/` | 78 | Local path vs URL abstraction |
| `data/config.go` | ~60 | JSON config load/save |
| `data/history.go` | ~80 | Browsing history with back/forward |

### Tier 2: Extract with minor refactoring

| Component | Effort | Notes |
|-----------|--------|-------|
| Tree state machine | 2-3 hrs | Remove `tea.Cmd` return types, use plain events |
| Event types | 1 hr | Already UI-agnostic, just need own package |
| File operations | 1 hr | `copyFile`, `copyDir`, paste logic |

### Tier 3: Parallel implementations (share interface only)

| Component | Why |
|-----------|-----|
| Rendering | Obviously different (lipgloss vs Gio) |
| Keyboard handling | Input systems differ, but map to same actions |
| Dialogs | Visual components, but same confirm/input patterns |

## Proposed Structure

```
github.com/jaycee1285/mdkit/
├── location/
│   └── location.go       # FromPath, FromURL, IsLocal, etc.
├── tree/
│   ├── state.go          # TreeState, Node, cursor logic
│   ├── events.go         # FileSelected, NewFileRequested, etc.
│   └── operations.go     # Expand, collapse, navigation
├── fileops/
│   └── fileops.go        # Copy, move, delete with error handling
└── history/
    └── history.go        # Back/forward navigation stack
```

Both apps import `mdkit`. crustdown wraps with BubbleTea, giopad wraps with Gio.

## Effort Estimate

| Task | Time |
|------|------|
| Extract mdkit from crustdown | 4-6 hrs |
| Update crustdown to import mdkit | 2-3 hrs |
| Update giopad to import mdkit | 3-4 hrs |
| **Total** | ~12 hrs |

## Risks / Honest Assessment

1. **Two consumers = coordination tax.** Changes to mdkit need testing in both apps. Acceptable given they're both yours.

2. **Premature abstraction?** giopad is younger. Some of its patterns might diverge. Mitigation: start with Tier 1 (location, history) which are stable.

3. **Go module complexity.** Three repos instead of two. Could alternatively vendor mdkit into both, but loses the "fix once" benefit.

## Recommendation

Start small:

1. **Now:** Copy `location/` into giopad as `internal/location`. Zero risk, immediate benefit.

2. **Next session:** Extract tree state machine. This is where the real value is - crustdown's tree handles refresh-while-preserving-expanded-state, cursor bounds, clipboard operations.

3. **Later:** If both apps are actively developed, promote to shared module.

Don't over-engineer the shared infrastructure until giopad has enough features to validate the abstraction boundaries.

---

*Reviewed: 2026-02-09*
*Source insight: Product review of crustdown/giopad overlap*
