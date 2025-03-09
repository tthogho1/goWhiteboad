package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
)

type whiteboard struct {
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

func newWhiteboard() *whiteboard {
	w := &whiteboard{
		lines: []line{},
		currentLine: line{
			points: []fyne.Position{},
			color:  color.Black,
			width:  2,
		},
		objects: []fyne.CanvasObject{},
	}
	w.ExtendBaseWidget(w)
	return w
}

// Cursorを実装することで、MouseMovedイベントを有効にする
func (w *whiteboard) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

func (w *whiteboard) SaveAsPNG(filename string, width, height int) error {
	dc := gg.NewContext(width, height)
	dc.SetColor(color.White)
	dc.Clear()

	for _, l := range w.lines {
		dc.SetColor(l.color)
		dc.SetLineWidth(float64(l.width))

		if len(l.points) > 1 {
			dc.MoveTo(float64(l.points[0].X), float64(l.points[0].Y))
			for _, p := range l.points[1:] {
				dc.LineTo(float64(p.X), float64(p.Y))
			}
			dc.Stroke()
		}
	}

	for _, obj := range w.objects {
		// オブジェクトの描画ロジックをここに追加
		// 例: 四角形、円など
		println(obj)
	}

	return dc.SavePNG(filename)
}

func (w *whiteboard) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(color.White)
	renderer := &whiteboardRenderer{
		whiteboard: w,
		background: background,
		objects:    []fyne.CanvasObject{background},
	}
	return renderer
}

// desktop.Hoverable を実装
func (w *whiteboard) MouseIn(*desktop.MouseEvent) {
	fmt.Fprintln(os.Stdout, "mouse in")
}

func (w *whiteboard) MouseOut() {
	fmt.Fprintln(os.Stdout, "mouse out")
}

// desktop.Mousable を実装
func (w *whiteboard) MouseDown(ev *desktop.MouseEvent) {
	w.drawing = true
	w.currentLine = line{
		points: []fyne.Position{ev.Position},
		color:  color.Black,
		width:  2,
	}
	fmt.Fprintf(os.Stdout, "mousedown at %v\n", ev.Position)
	w.Refresh()
}

func (w *whiteboard) MouseUp(ev *desktop.MouseEvent) {
	w.drawing = false
	if len(w.currentLine.points) > 0 {
		w.lines = append(w.lines, w.currentLine)
	}
	w.currentLine = line{
		points: []fyne.Position{},
		color:  color.Black,
		width:  2,
	}
	fmt.Fprintf(os.Stdout, "mouseup at %v\n", ev.Position)
	w.Refresh()
}

func (w *whiteboard) MouseMoved(ev *desktop.MouseEvent) {
	fmt.Fprintf(os.Stdout, "mousemove at %v\n", ev.Position)
	if w.drawing {
		w.currentLine.points = append(w.currentLine.points, ev.Position)
		w.Refresh()
	}
}

// Tappable インターフェースの実装
func (w *whiteboard) Tapped(ev *fyne.PointEvent) {
	fmt.Fprintf(os.Stdout, "tapped at %v\n", ev.Position)
}

type whiteboardRenderer struct {
	whiteboard *whiteboard
	background *canvas.Rectangle
	lines      []fyne.CanvasObject
	objects    []fyne.CanvasObject
}

func (r *whiteboardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 300)
}

func (r *whiteboardRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
}

func (r *whiteboardRenderer) Refresh() {
	r.lines = nil

	// 既存の線を描画
	for _, l := range r.whiteboard.lines {
		objects := createLineObjects(l)
		r.lines = append(r.lines, objects...)
	}

	// 現在描画中の線を描画
	if len(r.whiteboard.currentLine.points) > 0 {
		objects := createLineObjects(r.whiteboard.currentLine)
		r.lines = append(r.lines, objects...)
	}

	// objects リストを更新
	r.objects = append([]fyne.CanvasObject{r.background}, r.lines...)

	canvas.Refresh(r.whiteboard)
}

func (r *whiteboardRenderer) BackgroundColor() color.Color {
	return color.White
}

func (r *whiteboardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *whiteboardRenderer) Destroy() {
}

func createLineObjects(l line) []fyne.CanvasObject {
	var objects []fyne.CanvasObject

	if len(l.points) < 2 {
		return objects
	}

	for i := 1; i < len(l.points); i++ {
		p1 := l.points[i-1]
		p2 := l.points[i]

		line := canvas.NewLine(l.color)
		line.StrokeWidth = l.width
		line.Position1 = p1
		line.Position2 = p2

		// 線分の位置とサイズを設定
		x := fyne.Min(p1.X, p2.X)
		y := fyne.Min(p1.Y, p2.Y)
		width := float64(math.Abs(float64(p2.X - p1.X)))
		height := float64(math.Abs(float64(p2.Y - p1.Y)))

		// 線が垂直または水平の場合の調整
		if width == 0 {
			width = float64(line.StrokeWidth)
		}
		if height == 0 {
			height = float64(line.StrokeWidth)
		}

		line.Move(fyne.NewPos(x, y))
		line.Resize(fyne.NewSize(float32(width), float32(height)))

		objects = append(objects, line)
	}

	return objects
}
