package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Create form with settings
func ShowSettingDialog(w fyne.Window, board *whiteboard) {

	penColorSelect := widget.NewSelect([]string{"Black", "Red", "Blue", "Green"}, nil)
	penColorSelect.SetSelected("Black")

	// default width is set to 2
	DEFAULT_PEN_WIDTH := 2.0
	penWidthSlider := widget.NewSlider(1, 10)
	penWidthSlider.SetValue(DEFAULT_PEN_WIDTH) // Default width
	penWidthLabel := widget.NewLabel(fmt.Sprintf("%.0f", DEFAULT_PEN_WIDTH))

	penWidthSlider.OnChanged = func(value float64) {
    penWidthLabel.SetText(fmt.Sprintf("%.0f", value))
  }

	var customDialog dialog.Dialog

	// Create buttons for input forms
	EditSystemPrompt := widget.NewButton("System Prompt", func() {
		showInputForm1(w, board)
	})

	EditUserPrompt := widget.NewButton("User Prompt", func() {
		showInputForm2(w, board)
	})

	// Create button container with horizontal layout
	buttonContainer := container.New(layout.NewHBoxLayout(), EditSystemPrompt, EditUserPrompt)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Pen Color", Widget: penColorSelect},
			{Text: "Pen Width", Widget: container.NewBorder(nil, nil, nil, penWidthLabel, penWidthSlider)},
			{Text: "Additional Options", Widget: buttonContainer},
		},
		OnSubmit: func() {
			// Apply settings to the whiteboard
			var penColor color.Color
			switch penColorSelect.Selected {
			case "Red":
				penColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			case "Blue":
				penColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
			case "Green":
				penColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
			default:
				penColor = color.Black
			}

			// Update the current line settings
			board.currentLine.color = penColor
			board.currentLine.width = float32(penWidthSlider.Value)

			// Close the dialog
			if customDialog != nil {
				customDialog.Hide()
			}
			w.Canvas().Refresh(board)
		},
		OnCancel: func() {
			// Close the dialog
			if customDialog != nil {
				customDialog.Hide()
			}
		},
	}

	// Create and show the dialog
	customDialog = dialog.NewCustomWithoutButtons("Settings", form, w)
	customDialog.Resize(fyne.NewSize(300, 250))
	customDialog.Show()
}

// Input form for the first button
func showInputForm1(w fyne.Window, board *whiteboard) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter custom setting 1")

	var inputDialog dialog.Dialog

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Setting 1", Widget: entry},
		},
		OnSubmit: func() {
			// Process the input
			_ = entry.Text
			// Here you can add code to handle the input value
			// For example: board.customSetting1 = entry.Text

			if inputDialog != nil {
				inputDialog.Hide()
			}
		},
		OnCancel: func() {
			if inputDialog != nil {
				inputDialog.Hide()
			}
		},
	}

	inputDialog = dialog.NewCustom("System Prompt", "Close", container.NewVBox(form), w)

	inputDialog.Resize(fyne.NewSize(300, 150))
	inputDialog.Show()
}

// Input form for the second button
func showInputForm2(w fyne.Window, board *whiteboard) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter custom setting 2")

	var inputDialog dialog.Dialog

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Setting 2", Widget: entry},
		},
		OnSubmit: func() {
			// Process the input
			_ = entry.Text
			// Here you can add code to handle the input value
			// For example: board.customSetting2 = entry.Text

			if inputDialog != nil {
				inputDialog.Hide()
			}
		},
		OnCancel: func() {
			if inputDialog != nil {
				inputDialog.Hide()
			}
		},
	}

	// カスタムウィジェットを作成
	content := container.NewVBox(
		widget.NewLabel("User Prompt"),
		form,
	)

	// モーダルをカスタムで表示
	modal := widget.NewModalPopUp(content, w.Canvas())
	modal.Resize(fyne.NewSize(300, 150))
	modal.Show()

}
