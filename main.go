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

	// 画像生成ボタン
	saveButton := widget.NewButton("SavePng", func() {
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
		imagePath := "whiteboard.png" // 読み込むPNG画像のファイルパスを指定
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			log.Fatalf("画像の読み込みに失敗しました: %v", err)
		}

		htmlContent := util.SendImage(imageData)
		// TODO: HTMLを画面に表示する
		w := webview.New(true)
		w.SetTitle("Whiteboard")
		w.SetSize(int(BOARD_WIDTH), int(BOARD_HEIGHT), webview.HintNone)
		w.SetHtml(htmlContent)
		w.Run()
	})

	// メインコンテナ
	content := container.NewBorder(
		container.NewHBox(clearButton, saveButton, backButton, sendButton),
		nil,
		nil,
		nil,
		centerContainer,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(BOARD_WIDTH+100, BOARD_HEIGHT+100))
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyEscape {
			a.Quit()
		}
	})

	w.ShowAndRun()
}
