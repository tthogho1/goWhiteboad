package main

import (
	"fmt"
	"goWhiteBoard/util"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
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

	// 初期描画オブジェクトを生成するために一度 Refresh を呼び出す
	board.Refresh()

	// これでインターフェース実装を確認
	_, isMousable := interface{}(board).(desktop.Mouseable)
	_, isHoverable := interface{}(board).(desktop.Hoverable)
	_, isCursable := interface{}(board).(desktop.Cursorable)

	fmt.Printf("Mousable: %v, Hoverable: %v, Cursable: %v\n", isMousable, isHoverable, isCursable)

	// Ensure the whiteboard can receive mouse events
	if !isMousable {
		log.Fatal("Whiteboard does not implement desktop.Mouseable")
	}

	// 画像表示用のコンテナ
	imageContainer := container.NewStack()

	// 現在表示中のコンテンツ（ボードまたは画像）
	var currentContent fyne.CanvasObject = board

	// ヘッダーの作成
	headerLabel := widget.NewLabel("Whiteboard App")
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}

	// ヘッダーのスタイル設定
	headerBg := canvas.NewRectangle(color.RGBA{230, 230, 230, 255})

	// ヘッダーコンテナをスタックに配置して背景色を適用
	headerStack := container.NewStack(
		headerBg,
	)

	// ヘッダーコンテナ全体（背景色付きヘッダーと区切り線）
	headerContainer := container.NewVBox(
		headerStack,
		widget.NewSeparator(), // ヘッダーと本文の区切り線
	)

	// コンテンツを更新する関数
	updateContent := func() {
		content := container.NewBorder(
			headerContainer,
			nil,
			nil,
			nil,
			currentContent,
		)
		w.SetContent(content)
	}

	// ツールバー
	clearButton := widget.NewButton("Clear", func() {
		board.lines = []line{}
		board.Refresh()

		// Update dimensions
		size := w.Canvas().Size()
		BOARD_WIDTH = size.Width
		BOARD_HEIGHT = size.Height
		fmt.Printf("Window size: %v x %v\n", BOARD_WIDTH, BOARD_HEIGHT)
	})

	// 画像生成ボタン
	saveButton := widget.NewButton("SavePng", func() {
		// Update dimensions before saving
		size := w.Canvas().Size()
		BOARD_WIDTH = size.Width
		BOARD_HEIGHT = size.Height

		board.SaveAsPNG("whiteboard.png", int(BOARD_WIDTH), int(BOARD_HEIGHT))

		img := canvas.NewImageFromFile("whiteboard.png")
		img.FillMode = canvas.ImageFillOriginal

		// 画像をコンテナに設定
		imageContainer.Objects = []fyne.CanvasObject{img}

		// メインコンテンツを画像に切り替え
		currentContent = imageContainer
		updateContent()
	})

	backButton := widget.NewButton("Back to Drawing", func() {
		// メインコンテンツをボードに切り替え
		currentContent = board
		updateContent()
	})

	// 画像送信ボタン。送信した結果をhtmlで受け取る。
	sendButton := widget.NewButton("Send", func() {
		// Update dimensions before sending
		size := w.Canvas().Size()
		BOARD_WIDTH = size.Width
		BOARD_HEIGHT = size.Height

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

	// ボタンコンテナ
	buttonContainer := container.NewHBox(
		clearButton,
		saveButton,
		backButton,
		sendButton,
		settingsButton,
	)

	// ヘッダーコンテンツ
	headerContent := container.NewHBox(
		headerLabel,
		layout.NewSpacer(), // 左右の間隔を空ける
		buttonContainer,
	)

	// ヘッダースタックにコンテンツを追加
	headerStack.Add(container.NewPadded(headerContent))

	// 初期コンテンツを設定
	content := container.NewBorder(
		headerContainer,
		nil,
		nil,
		nil,
		currentContent,
	)

	// 初期サイズを設定
	var lastSize fyne.Size

	// ウィンドウサイズ変更時のコールバック
	w.Canvas().SetOnTypedRune(func(r rune) {
		// This is just a hook to get called regularly
		currentSize := w.Canvas().Size()
		if currentSize != lastSize {
			lastSize = currentSize
			BOARD_WIDTH = currentSize.Width
			BOARD_HEIGHT = currentSize.Height
			fmt.Printf("Window size: %v x %v\n", BOARD_WIDTH, BOARD_HEIGHT)
		}
	})

	// ESCキーでアプリを終了
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Name == fyne.KeyEscape {
			a.Quit()
		}
	})
	w.SetContent(content)
	w.Resize(fyne.NewSize(BOARD_WIDTH+100, BOARD_HEIGHT+100))

	// アプリを実行
	w.ShowAndRun()
}
