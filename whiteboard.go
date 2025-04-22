package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Point represents a point on the whiteboard
type Point struct {
	X, Y float32
}

// Line represents a line on the whiteboard
type line struct {
	points []Point
	color  color.Color
	width  float32
}

// Whiteboard is a custom widget for drawing
type whiteboard struct {
	widget.BaseWidget
	lines       []line
	currentLine line
	drawing     bool
	lineColor   color.Color
	lineWidth   float32
	mutex       sync.Mutex // 複数のゴルーチンからのアクセスを保護
}

// NewWhiteboard creates a new whiteboard widget
func newWhiteboard() *whiteboard {
	w := &whiteboard{
		lines:     []line{},
		lineColor: color.RGBA{0, 0, 0, 255}, // Default: Black
		lineWidth: 2.0,                      // Default width
	}
	w.ExtendBaseWidget(w)
	return w
}

// Ensure the whiteboard implements the necessary interfaces
var _ desktop.Mouseable = (*whiteboard)(nil)
var _ fyne.Widget = (*whiteboard)(nil)
var _ desktop.Hoverable = (*whiteboard)(nil)
var _ desktop.Cursorable = (*whiteboard)(nil)

// CreateRenderer implements the fyne.Widget interface
func (w *whiteboard) CreateRenderer() fyne.WidgetRenderer {
	return &whiteboardRenderer{whiteboard: w}
}

// MouseDown implements desktop.Mouseable
func (w *whiteboard) MouseDown(ev *desktop.MouseEvent) {
	w.drawing = true
	w.currentLine = line{
		points: []Point{{X: ev.Position.X, Y: ev.Position.Y}},
		color:  w.lineColor,
		width:  w.lineWidth,
	}
}

// MouseUp implements desktop.Mouseable
func (w *whiteboard) MouseUp(ev *desktop.MouseEvent) {
	if w.drawing {
		w.drawing = false
		w.lines = append(w.lines, w.currentLine)
		w.currentLine = line{}
	}
}

// MouseMoved implements desktop.Mouseable
func (w *whiteboard) MouseMoved(ev *desktop.MouseEvent) {
	if w.drawing {
		w.currentLine.points = append(w.currentLine.points, Point{X: ev.Position.X, Y: ev.Position.Y})
		w.Refresh()
	}
}

// MouseIn implements desktop.Hoverable
func (w *whiteboard) MouseIn(*desktop.MouseEvent) {
}

// MouseOut implements desktop.Hoverable
func (w *whiteboard) MouseOut() {
}

// Cursor implements desktop.Cursorable
func (w *whiteboard) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

// SetLineColor sets the color for new lines
func (w *whiteboard) SetLineColor(c color.Color) {
	w.lineColor = c
}

// SetLineWidth sets the width for new lines
func (w *whiteboard) SetLineWidth(width float32) {
	w.lineWidth = width
}

// SaveAsPNG saves the whiteboard as a PNG image
func (w *whiteboard) SaveAsPNG(filename string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with white background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}

	// Draw all lines
	tempLines := make([]line, len(w.lines))
	copy(tempLines, w.lines)
	current := w.currentLine
	drawing := w.drawing

	for _, l := range tempLines {
		drawLine(img, l)
	}

	// Draw current line if drawing
	if drawing {
		drawLine(img, current)
	}

	// Save to file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// drawLine draws a line on the image
func drawLine(img *image.RGBA, l line) {
	if len(l.points) < 2 {
		return
	}

	// Convert color.Color to RGBA
	r, g, b, a := l.color.RGBA()
	lineColor := color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}

	// Draw line segments
	for i := 1; i < len(l.points); i++ {
		p1 := l.points[i-1]
		p2 := l.points[i]
		drawLineSegment(img, int(p1.X), int(p1.Y), int(p2.X), int(p2.Y), lineColor, int(l.width))
	}
}

// drawLineSegment draws a line segment using Bresenham's algorithm with thickness
func drawLineSegment(img *image.RGBA, x0, y0, x1, y1 int, col color.RGBA, thickness int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	bounds := img.Bounds()
	radius := thickness / 2
	if radius < 1 {
		radius = 1
	}

	for {
		for y := -radius; y <= radius; y++ {
			for x := -radius; x <= radius; x++ {
				if x*x+y*y <= radius*radius {
					px, py := x0+x, y0+y
					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						img.Set(px, py, col)
					}
				}
			}
		}

		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// whiteboardRenderer implements the fyne.WidgetRenderer interface
type whiteboardRenderer struct {
	whiteboard *whiteboard
	objects    []fyne.CanvasObject // 描画するオブジェクトをキャッシュ
}

// MinSize implements fyne.WidgetRenderer
func (r *whiteboardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 200)
}

// Layout implements fyne.WidgetRenderer
func (r *whiteboardRenderer) Layout(size fyne.Size) {
	// サイズが変わった場合にオブジェクトの再配置などを行う場合はここに記述
}

// Refresh implements fyne.WidgetRenderer
func (r *whiteboardRenderer) Refresh() {
	fmt.Println("whiteboardRenderer.Refresh() called") // ログを追加
	r.updateObjects()                                  // Refresh 時に描画オブジェクトを更新
	canvas.Refresh(r.whiteboard)
	fmt.Println("whiteboardRenderer.Refresh() finished")
}

// Objects implements fyne.WidgetRenderer
func (r *whiteboardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy implements fyne.WidgetRenderer
func (r *whiteboardRenderer) Destroy() {
	// 特にクリーンアップ処理はなし
}

// updateObjects は whiteboard の線のデータを canvas.Line オブジェクトに変換してキャッシュする
func (r *whiteboardRenderer) updateObjects() {
	fmt.Println("updateObjects()")
	r.whiteboard.mutex.Lock()
	defer r.whiteboard.mutex.Unlock()

	r.objects = make([]fyne.CanvasObject, 0, len(r.whiteboard.lines)+1) // 描画オブジェクトのスライスを初期化

	// 描画済みの線を canvas.Line オブジェクトに変換
	for _, l := range r.whiteboard.lines {
		if len(l.points) >= 2 {
			for i := 0; i < len(l.points)-1; i++ {
				line := canvas.NewLine(l.color)
				line.StrokeWidth = l.width
				line.Position1 = fyne.NewPos(l.points[i].X, l.points[i].Y)
				line.Position2 = fyne.NewPos(l.points[i+1].X, l.points[i+1].Y)
				r.objects = append(r.objects, line)
			}
		}
	}

	// 現在描画中の線も追加
	if r.whiteboard.drawing && len(r.whiteboard.currentLine.points) >= 2 {
		for i := 0; i < len(r.whiteboard.currentLine.points)-1; i++ {
			currentLine := canvas.NewLine(r.whiteboard.currentLine.color)
			currentLine.StrokeWidth = r.whiteboard.currentLine.width
			currentLine.Position1 = fyne.NewPos(r.whiteboard.currentLine.points[i].X, r.whiteboard.currentLine.points[i].Y)
			currentLine.Position2 = fyne.NewPos(r.whiteboard.currentLine.points[i+1].X, r.whiteboard.currentLine.points[i+1].Y)
			r.objects = append(r.objects, currentLine)
		}
	}
}
