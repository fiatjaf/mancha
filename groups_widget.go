package main

import (
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type GroupsWidget struct {
	widget        *widget.List
	overlay       *fyne.Container
	searchResults *widget.List
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
			img := canvas.NewImageFromImage(neutralImage)
			img.SetMinSize(fyne.NewSize(36, 36))

			return container.NewHBox(
				img,
				widget.NewLabel("template"),
				widget.NewLabel("template"),
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
				Bold:   true,
				Italic: false,
			}

			container.Objects[2].(*widget.Label).SetText(group.ID)
			container.Objects[2].(*widget.Label).TextStyle = fyne.TextStyle{
				Bold:      false,
				Italic:    true,
				Monospace: true,
			}
			container.Objects[2].(*widget.Label).Alignment = fyne.TextAlignTrailing
			container.Objects[2].(*widget.Label).Importance = widget.LowImportance
		},
	)

	gw.widget.OnSelected = func(i widget.ListItemID) {
		state.selected = state.groups[i]
		getInputWidget().enable()
		getMessagesWidget().widget.Refresh()
		getMessagesWidget().widget.ScrollToBottom()
	}

	gw.overlay = container.NewCenter(
		widget.NewButtonWithIcon("Join Group", theme.ContentAddIcon(), showAddGroupDialog),
	)

	return gw
}

func showAddGroupDialog() {
	urlEntry := newEnhancedEntry()
	urlEntry.PlaceHolder = "groups.nostr.com"

	naddrEntry := newEnhancedEntry()
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

				getGroupsWidget().widget.Refresh()
				getGroupsWidget().overlay.Hide()

			}
		}

		if urlEntry.Text != "" {
			events := pool.SubManyEose(context.Background(), []string{urlEntry.Text}, nostr.Filters{
				{Limit: 40, Kinds: []int{nostr.KindSimpleGroupMetadata}},
			})
			groups := make([]*Group, 0, 40)
			for ie := range events {
				group := &Group{
					Relay: urlEntry.Text,
				}
				group.Group.MergeInMetadataEvent(ie.Event)
				groups = append(groups, group)
			}
			fmt.Println(groups)
		}
	}, w)
}
