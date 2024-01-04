package main

import (
	"encoding/hex"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasb-eyer/go-colorful"
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
			name := canvas.NewText("template", color.Transparent)
			name.TextStyle.Bold = true
			name.Alignment = fyne.TextAlignTrailing

			message := widget.NewLabel("template")
			message.Alignment = fyne.TextAlignLeading
			message.Wrapping = fyne.TextWrapWord

			border := container.NewBorder(nil, nil, name, nil, message)

			return border
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			evt := state.selected.messages[i]
			container := o.(*fyne.Container)

			name := "  [ " + evt.PubKey[0:8] + " ]"
			if metadata, _ := people.Load(evt.PubKey); metadata != nil {
				name = "  " + metadata.ShortName()
			}
			mw.widget.SetItemHeight(i, o.(*fyne.Container).Objects[0].(*widget.Label).MinSize().Height)

			nameLabel := container.Objects[1].(*canvas.Text)
			nameLabel.Text = name
			lastByte, _ := hex.DecodeString(evt.PubKey[62:])
			nameLabel.Color = colorful.Hsl(float64(lastByte[0])/256*360, 0.41, 0.55)

			container.Objects[0].(*widget.Label).Text = evt.Content
		},
	)

	return mw
}
