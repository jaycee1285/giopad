package toggle

import (
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

// ThemeToggle is a sun/moon button to switch themes
type ThemeToggle struct {
	click widget.Clickable
}

// New creates a new ThemeToggle
func New() *ThemeToggle {
	return &ThemeToggle{}
}

// Layout renders the toggle button
func (t *ThemeToggle) Layout(gtx C, th *material.Theme) D {
	// Handle click
	if t.click.Clicked(gtx) {
		app.ToggleTheme()
	}

	return t.click.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
			Left:   unit.Dp(12),
			Right:  unit.Dp(12),
		}.Layout(gtx, func(gtx C) D {
			icon := "☀"  // Sun for "switch to light"
			if !app.CurrentTheme.IsDark {
				icon = "🌙" // Moon for "switch to dark"
			}
			label := material.Body1(th, icon)
			label.Color = app.Comment()
			return label.Layout(gtx)
		})
	})
}
