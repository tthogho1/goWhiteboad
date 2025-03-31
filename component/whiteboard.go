package component

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Whiteboard struct {
	widget.BaseWidget
	lines       []line
	currentLine line
	drawing     bool
	objects     []fyne.CanvasObject
}

type line struct {
	points []fyne.Position
	color  color.Color
	width  float32
}
