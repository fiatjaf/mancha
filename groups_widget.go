package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GroupsWidget struct {
	widget *widget.List
}

func makeGroupsWidget() *GroupsWidget {
	gw := &GroupsWidget{}

	gw.widget = widget.NewList(
		func() int {
			return len(state.groups)
		},
		func() fyne.CanvasObject {
			btn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
				fmt.Println("button clicked")
			})
			btn.Importance = widget.LowImportance

			img := canvas.NewImageFromImage(neutralImage)
			img.SetMinSize(fyne.NewSize(btn.MinSize().Height, btn.MinSize().Height))

			return container.NewHBox(
				img,
				widget.NewLabel("template"),
				layout.NewSpacer(),
				btn,
			)
		},
		func(lii widget.ListItemID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)

			group := state.groups[lii]
			if group.Picture != "" {
				container.Objects[0].(*canvas.Image).Image = imageFromURL(group.Picture)
			}

			container.Objects[1].(*widget.Label).SetText(group.Name)
			container.Objects[1].(*widget.Label).TextStyle = fyne.TextStyle{
				Bold:   false,
				Italic: false,
			}

			container.Objects[2].Hide()
		},
	)

	gw.widget.OnSelected = func(i widget.ListItemID) {
		state.selected = state.groups[i]
		messagesWidget.widget.Refresh()
		messagesWidget.widget.ScrollToBottom()
	}

	return gw
}
