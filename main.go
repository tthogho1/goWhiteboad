package main

import (
	"fmt"
	"goWhiteBoard/util"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/joho/godotenv"
	webview "github.com/webview/webview_go"
)

// Initial board dimensions that will be updated when window is resized
var (
	BOARD_WIDTH  float32 = 800
	BOARD_HEIGHT float32 = 600
)

// 修正：完全な実装での確認コード
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	a := app.New()
	w := a.NewWindow("Whiteboard")
	w.Resize(fyne.NewSize(BOARD_WIDTH, BOARD_HEIGHT))

	board := newWhiteboard()

	// これでインターフェース実装を確認
	_, isMousable := interface{}(board).(desktop.Mouseable)
	_, isHoverable := interface{}(board).(desktop.Hoverable)
	_, isCursable := interface{}(board).(desktop.Cursorable)

	fmt.Printf("Mousable: %v, Hoverable: %v, Cursable: %v\n", isMousable, isHoverable, isCursable)

	//
	centerContainer := container.NewStack(board)

	// ツールバー
	clearButton := widget.NewButton("Clear", func() {
		board.lines = []line{}
		board.Refresh()
	})

	// Create a function to update board dimensions based on window size
	updateBoardDimensions := func() {
		size := w.Canvas().Size()
		BOARD_WIDTH = size.Width
		BOARD_HEIGHT = size.Height
		fmt.Printf("Window size: %v x %v\n", BOARD_WIDTH, BOARD_HEIGHT)
	}

	// Add a resize listener using a custom approach
	resizeListener := widget.NewLabel("")
	resizeListener.Hide()
	resizeListener.Resize(fyne.NewSize(1, 1))

	// Add a resize callback to the window content
	var lastSize fyne.Size

	// Create a custom widget that will detect size changes
	sizeDetector := widget.NewLabel("")
	sizeDetector.Hide()

	// Set up a callback that will be triggered on each render
	w.Canvas().SetOnTypedRune(func(r rune) {
		// This is just a hook to get called regularly
		currentSize := w.Canvas().Size()
		if currentSize != lastSize {
			lastSize = currentSize
			updateBoardDimensions()
		}
	})

	// Also update dimensions when buttons are clicked
	clearButton.OnTapped = func() {
		updateBoardDimensions()
		board.lines = []line{}
		board.Refresh()
	}

	// 画像生成ボタン
	saveButton := widget.NewButton("SavePng", func() {
		// Update dimensions before saving
		updateBoardDimensions()
		board.SaveAsPNG("whiteboard.png", int(BOARD_WIDTH), int(BOARD_HEIGHT))

		img := canvas.NewImageFromFile("whiteboard.png")
		img.FillMode = canvas.ImageFillOriginal

		// 中央部分を画像に置き換える
		centerContainer.Objects = []fyne.CanvasObject{img}
		centerContainer.Refresh()
	})

	backButton := widget.NewButton("Back to Drawing", func() {
		centerContainer.Objects = []fyne.CanvasObject{board}
		centerContainer.Refresh()
	})

	// 画像送信ボタン。送信した結果をhtmlで受け取る。
	sendButton := widget.NewButton("Send", func() {
		// Update dimensions before sending
		updateBoardDimensions()

		imagePath := "whiteboard.png" // 読み込むPNG画像のファイルパスを指定
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			log.Fatalf("画像の読み込みに失敗しました: %v", err)
		}

		htmlContent := util.SendImage(imageData)
		// TODO: HTMLを画面に表示する
		wv := webview.New(true)
		wv.SetTitle("Whiteboard")
		wv.SetSize(int(BOARD_WIDTH), int(BOARD_HEIGHT), webview.HintNone)
		wv.SetHtml(htmlContent)
		wv.Run()
	})

	// 設定ボタン
	settingsButton := widget.NewButton("Settings", func() {
		ShowSettingDialog(w, board)
	})

	// メインコンテナ
	content := container.NewBorder(
		container.NewHBox(clearButton, saveButton, backButton, sendButton, settingsButton),
		nil,
		nil,
		nil,
		container.NewStack(centerContainer, sizeDetector, resizeListener),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(BOARD_WIDTH+100, BOARD_HEIGHT+100))
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	// Initial update of dimensions
	updateBoardDimensions()

	w.ShowAndRun()
}
