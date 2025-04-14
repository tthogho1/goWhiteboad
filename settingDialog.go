package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Create form with settings
func ShowSettingDialog(w fyne.Window, board *whiteboard) {

	penColorSelect := widget.NewSelect([]string{"Black", "Red", "Blue", "Green"}, nil)
	penColorSelect.SetSelected("Black")

	penWidthSlider := widget.NewSlider(1, 10)
	penWidthSlider.SetValue(2) // Default width
	penWidthLabel := widget.NewLabel("2")

	var customDialog dialog.Dialog

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Pen Color", Widget: penColorSelect},
			{Text: "Pen Width", Widget: container.NewBorder(nil, nil, nil, penWidthLabel, penWidthSlider)},
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
	customDialog.Resize(fyne.NewSize(300, 200))
	customDialog.Show()
}
