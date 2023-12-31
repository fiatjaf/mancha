package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeBottomWidget() *fyne.Container {
	// setup the right side of the window
	chatInputWidget := widget.NewMultiLineEntry()
	chatInputWidget.Wrapping = fyne.TextWrapWord
	chatInputWidget.PlaceHolder = "Your message here... shift+enter to Submit"
	chatInputWidget.OnSubmitted = func(s string) {
		go func() {
			if s == "" {
				return
			}
			chatInputWidget.SetText("")
			if err := publishChat(state.selected, s); err != nil {
				// TODO show a message to user about this error
				fmt.Println("failed to publish:", err)
			}
		}()
	}

	submitChatButtonWidget := widget.NewButton("Submit", func() {
		message := chatInputWidget.Text
		if message == "" {
			return
		}
		go func() {
			chatInputWidget.SetText("")
			publishChat(state.selected, message)
		}()
	})

	return container.NewBorder(nil, nil, nil, submitChatButtonWidget, chatInputWidget)
}
