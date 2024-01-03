package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MessagesWidget struct {
	widget *widget.List
}

var _messagesWidget *MessagesWidget

func getMessagesWidget() *MessagesWidget {
	if _messagesWidget == nil {
		_messagesWidget = makeMessagesWidget()
	}
	return _messagesWidget
}

func makeMessagesWidget() *MessagesWidget {
	mw := &MessagesWidget{}

	mw.widget = widget.NewList(
		func() int {
			if state.selected == nil {
				return 0
			}
			return len(state.selected.messages)
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
			evt := state.selected.messages[i]
			container := o.(*fyne.Container)

			if metadata, _ := people.Load(evt.PubKey); metadata != nil {
				container.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Text).Text = fmt.Sprintf("[ %s ]", metadata.ShortName())
			}
			container.Objects[0].(*widget.Label).SetText(evt.Content)
			mw.widget.SetItemHeight(i, o.(*fyne.Container).Objects[0].(*widget.Label).MinSize().Height)
		},
	)

	return mw
}
