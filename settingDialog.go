package main

import (
	"fmt"
	"goWhiteBoard/config"
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
		showSystemPromptForm(w, board)
	})

	EditUserPrompt := widget.NewButton("User Prompt", func() {
		showUserPromptForm(w, board)
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
func showSystemPromptForm(w fyne.Window, board *whiteboard) {
	systemEntry := widget.NewMultiLineEntry()
	systemEntry.SetText(config.APISystemMessage) // Set default value

	var systemInputModal *widget.PopUp	

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "SystemPrompt", Widget: systemEntry},
		},
		OnSubmit: func() {
			// Process the input
			config.APISystemMessage = systemEntry.Text
			// Here you can add code to handle the input value
			// For example: board.customSetting1 = entry.Text

			if systemInputModal != nil {
				systemInputModal.Hide()
			}
		},
		OnCancel: func() {
			if systemInputModal != nil {
				systemInputModal.Hide()
			}
		},
	}

		// カスタムウィジェットを作成
	content := container.NewVBox(
		widget.NewLabel("System Prompt"),
		form,
	)

	systemInputModal = widget.NewModalPopUp(content, w.Canvas())
	systemInputModal.Resize(fyne.NewSize(600, 200))
	systemInputModal.Show()
}


// Input form for the second button
func showUserPromptForm(w fyne.Window, board *whiteboard) {
	userEntry := widget.NewMultiLineEntry()
	userEntry.SetText(config.APIUserMessage) // Set default value

	var userInputModal *widget.PopUp	

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "UserPrompt", Widget: userEntry},
		},
		OnSubmit: func() {
			// Process the input
			config.APIUserMessage = userEntry.Text
			// Here you can add code to handle the input value
			// For example: board.customSetting2 = entry.Text

			if userInputModal != nil {
				userInputModal.Hide()
			}
		},
		OnCancel: func() {
			if userInputModal != nil {
				userInputModal.Hide()
			}
		},
	}

	// カスタムウィジェットを作成
	content := container.NewVBox(
		widget.NewLabel("User Prompt"),
		form,
	)

	// モーダルをカスタムで表示
	userInputModal = widget.NewModalPopUp(content, w.Canvas())
	userInputModal.Resize(fyne.NewSize(600, 200))
	userInputModal.Show()

}
