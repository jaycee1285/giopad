package app

import (
	"image/color"

	"gioui.org/widget/material"
)

// Theme holds the current color palette
type Theme struct {
	Background color.NRGBA
	Surface    color.NRGBA
	Foreground color.NRGBA
	Comment    color.NRGBA
	Selection  color.NRGBA
	Red        color.NRGBA
	Green      color.NRGBA
	Yellow     color.NRGBA
	Blue       color.NRGBA
	Purple     color.NRGBA
	Cyan       color.NRGBA
	IsDark     bool
}

// Ayu Mirage (dark)
var AyuMirage = Theme{
	Background: rgb(0x1f2430),
	Surface:    rgb(0x232936),
	Foreground: rgb(0xcbccc6),
	Comment:    rgb(0x5c6773),
	Selection:  rgb(0x343f4c),
	Red:        rgb(0xf28b82),
	Green:      rgb(0xbae67e),
	Yellow:     rgb(0xffd580),
	Blue:       rgb(0x73d0ff),
	Purple:     rgb(0xd4bfff),
	Cyan:       rgb(0x95e6cb),
	IsDark:     true,
}

// Ayu Light
var AyuLight = Theme{
	Background: rgb(0xfafafa),
	Surface:    rgb(0xffffff),
	Foreground: rgb(0x575f66),
	Comment:    rgb(0xabb0b6),
	Selection:  rgb(0xd1e4f4),
	Red:        rgb(0xf07171),
	Green:      rgb(0x86b300),
	Yellow:     rgb(0xf2ae49),
	Blue:       rgb(0x399ee6),
	Purple:     rgb(0xa37acc),
	Cyan:       rgb(0x4cbf99),
	IsDark:     false,
}

// Current theme - mutable
var CurrentTheme = AyuMirage

// Color getters - return current theme colors
func Background() color.NRGBA  { return CurrentTheme.Background }
func Surface() color.NRGBA     { return CurrentTheme.Surface }
func Foreground() color.NRGBA  { return CurrentTheme.Foreground }
func Comment() color.NRGBA     { return CurrentTheme.Comment }
func Selection() color.NRGBA   { return CurrentTheme.Selection }
func Red() color.NRGBA         { return CurrentTheme.Red }
func Green() color.NRGBA       { return CurrentTheme.Green }
func Yellow() color.NRGBA      { return CurrentTheme.Yellow }
func Blue() color.NRGBA        { return CurrentTheme.Blue }
func Purple() color.NRGBA      { return CurrentTheme.Purple }
func Cyan() color.NRGBA        { return CurrentTheme.Cyan }
func Accent() color.NRGBA      { return CurrentTheme.Blue }

// ToggleTheme switches between light and dark
func ToggleTheme() {
	if CurrentTheme.IsDark {
		CurrentTheme = AyuLight
	} else {
		CurrentTheme = AyuMirage
	}
}

// AyuMirageTheme returns a material.Theme configured with current colors
func AyuMirageTheme() *material.Theme {
	th := material.NewTheme()
	th.Bg = CurrentTheme.Background
	th.Fg = CurrentTheme.Foreground
	th.ContrastBg = CurrentTheme.Blue
	th.ContrastFg = CurrentTheme.Background
	return th
}

func rgb(hex uint32) color.NRGBA {
	return color.NRGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: 0xff,
	}
}
