package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
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

// 修正：完全な実装での確認コード
func main() {
	a := app.New()
	w := a.NewWindow("Whiteboard")
	w.Resize(fyne.NewSize(800, 600))

	board := newWhiteboard()

	// これでインターフェース実装を確認
	_, isMousable := interface{}(board).(desktop.Mouseable)
	_, isHoverable := interface{}(board).(desktop.Hoverable)
	_, isCursable := interface{}(board).(desktop.Cursorable)

	fmt.Printf("Mousable: %v, Hoverable: %v, Cursable: %v\n", isMousable, isHoverable, isCursable)

	// ツールバー
	clearButton := widget.NewButton("Clear", func() {
		board.lines = []line{}
		board.Refresh()
	})

	// メインコンテナ
	content := container.NewBorder(
		container.NewHBox(clearButton),
		nil,
		nil,
		nil,
		board,
	)

	w.SetContent(content)
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	w.ShowAndRun()
}
