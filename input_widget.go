package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type InputWidget struct {
	widget *fyne.Container
	input  *middleClickHandlerEntry
	button *widget.Button
}

func makeInputWidget() *InputWidget {
	iw := &InputWidget{}

	iw.input = &middleClickHandlerEntry{}
	iw.input.ExtendBaseWidget(iw.input)
	iw.input.Wrapping = fyne.TextWrapOff
	iw.input.MultiLine = false
	iw.input.PlaceHolder = "Your message here... shift+enter to Submit"
	iw.input.OnSubmitted = func(s string) {
		go func() {
			if s == "" {
				return
			}
			iw.input.SetText("")
			if err := publishChat(state.selected, s); err != nil {
				// TODO show a message to user about this error
				fmt.Println("failed to publish:", err)
			}
		}()
	}

	iw.button = widget.NewButton("Submit", func() {
		message := iw.input.Text
		if message == "" {
			return
		}
		go func() {
			iw.input.SetText("")
			publishChat(state.selected, message)
		}()
	})

	iw.widget = container.NewBorder(nil, nil, nil, iw.button, iw.input)

	return iw
}

func (iw InputWidget) enable() {
	iw.input.Enable()
	iw.button.Enable()
}

func (iw InputWidget) onMount() {
	iw.input.Disable()
	iw.button.Disable()
}
