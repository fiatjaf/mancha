package main

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type GroupsWidget struct {
	widget  *widget.List
	overlay *fyne.Container
}

var _groupsWidget *GroupsWidget

func getGroupsWidget() *GroupsWidget {
	if _groupsWidget == nil {
		_groupsWidget = makeGroupsWidget()
	}
	return _groupsWidget
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
		getInputWidget().enable()
		getMessagesWidget().widget.Refresh()
		getMessagesWidget().widget.ScrollToBottom()
	}

	gw.overlay = container.NewCenter(
		widget.NewButtonWithIcon("Join Group", theme.StorageIcon(), func() {
			urlEntry := widget.NewEntry()
			urlEntry.PlaceHolder = "groups.nostr.com"

			naddrEntry := widget.NewEntry()
			naddrEntry.PlaceHolder = "naddr1..."

			dialog.ShowForm("Add group                                                                           ", "Add", "Cancel", []*widget.FormItem{ // Empty space Hack to make dialog bigger
				widget.NewFormItem("Relay URL:", urlEntry),
				widget.NewFormItem("or group code:", naddrEntry),
			}, func(b bool) {
				if !b {
					return
				}

				if naddrEntry.Text != "" {
					prefix, value, err := nip19.Decode(naddrEntry.Text)
					if err == nil && prefix == "naddr" {
						ent := value.(nostr.EntityPointer)
						if len(ent.Relays) != 1 {
							return
						}

						_, err := joinGroup(ent.Relays[0], ent.Identifier)
						if err != nil {
							log.Printf("error joining group: %s\n", err)
							return
						}

						gw.widget.Refresh()
						gw.overlay.Hide()

					}
				}

				if urlEntry.Text != "" {
					host := urlEntry.Text
					if strings.HasPrefix(host, "https://") {
						host = host[8:]
					}
					if strings.HasPrefix(host, "http://") {
						host = host[7:]
					}
					if strings.HasPrefix(host, "wss://") {
						host = host[6:]
					}
					if strings.HasPrefix(host, "ws://") {
						host = host[5:]
					}

					fmt.Println("relay", host)
				}
			}, w)
		}),
	)

	return gw
}
