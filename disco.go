package main

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// discoPalette defines the disco color rotation for alive cells.
var discoPalette = []color.Color{
	lipgloss.Color("#FF69B4"), // Hot Pink
	lipgloss.Color("#FFD700"), // Gold
	lipgloss.Color("#00BFFF"), // Electric Blue
	lipgloss.Color("#39FF14"), // Lime
	lipgloss.Color("#8B00FF"), // Violet
	lipgloss.Color("#FF6600"), // Orange
	lipgloss.Color("#00FFFF"), // Cyan
	lipgloss.Color("#FF1744"), // Red
}

const (
	aliveRune = "\u2588" // Full block character
	deadRune  = "\u00B7" // Middle dot
)

// taglines are rotating Bee Gees lyrics for the status bar.
var taglines = []string{
	"Ah, ha, ha, ha, stayin' alive",
	"Stayin' alive, stayin' alive",
	"Feel the city breakin' and everybody shakin'",
	"Life goin' nowhere, somebody help me",
	"Well, you can tell by the way I use my walk",
	"Whether you're a brother or whether you're a mother",
}

// CellStyle returns a lipgloss style for an alive cell based on generation and position.
func CellStyle(generation, x, y int) lipgloss.Style {
	idx := (generation + x + y) % len(discoPalette)
	return lipgloss.NewStyle().Foreground(discoPalette[idx]).Bold(true)
}

// DeadCellStyle returns a dim style for dead cells.
func DeadCellStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
}

// RenderCell renders a single cell as a styled string.
func RenderCell(alive bool, generation, x, y int) string {
	if alive {
		return CellStyle(generation, x, y).Render(aliveRune)
	}
	return DeadCellStyle().Render(deadRune)
}

// DiscoTagline returns a rotating tagline based on the generation number.
func DiscoTagline(generation int) string {
	return taglines[generation%len(taglines)]
}
