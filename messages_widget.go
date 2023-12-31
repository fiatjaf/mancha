package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MessagesWidget struct {
	widget *widget.List
}

func makeMessagesWidget() *MessagesWidget {
	mw := &MessagesWidget{}

	mw.widget = widget.NewList(
		func() int {
			return 0
		},
		func() fyne.CanvasObject {
			pubKey := canvas.NewText("template", color.RGBA{139, 190, 178, 255})
			pubKey.TextStyle.Bold = true
			pubKey.Alignment = fyne.TextAlignLeading

			message := widget.NewLabel("template")
			message.Alignment = fyne.TextAlignLeading
			message.Wrapping = fyne.TextWrapWord

			vbx := container.NewVBox(container.NewPadded(pubKey))
			border := container.NewBorder(nil, nil, vbx, nil, message)

			return border
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			chatMessage := state.selected.messages[i]

			var name string
			if metadata, _ := people.Load(chatMessage.PubKey); metadata != nil && metadata.Name != "" {
				name = fmt.Sprintf("[ %s ]", strings.TrimSpace(metadata.Name))
			} else {
				name = fmt.Sprintf("[ %s ]", chatMessage.PubKey[len(chatMessage.PubKey)-8:])
			}
			message := chatMessage.Content
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Text).Text = name
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(message)
			mw.widget.SetItemHeight(i, o.(*fyne.Container).Objects[0].(*widget.Label).MinSize().Height)
		},
	)

	return mw
}
