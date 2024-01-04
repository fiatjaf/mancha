package main

import (
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lucasb-eyer/go-colorful"
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
			img := canvas.NewImageFromImage(transparentBox)
			img.SetMinSize(fyne.NewSize(36, 36))

			idLabel := canvas.NewText("id", theme.DisabledColor())
			idLabel.Alignment = fyne.TextAlignTrailing
			idLabel.TextStyle = fyne.TextStyle{
				Bold:      false,
				Italic:    true,
				Monospace: true,
			}

			return container.NewHBox(
				img,
				widget.NewLabelWithStyle("template", fyne.TextAlignLeading, fyne.TextStyle{
					Bold:   true,
					Italic: false,
				}),
				layout.NewSpacer(),
				idLabel,
			)
		},
		func(lii widget.ListItemID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			group := state.groups[lii]
			if group.Picture == "" {
				byte := []byte(group.Name[0:1])[0]
				hue := float64(byte*byte*byte) / 256 * 360
				container.Objects[0].(*canvas.Image).Image = boxImage(colorful.Hsl(hue, 0.48, 0.52))
			} else {
				container.Objects[0].(*canvas.Image).Image = imageFromURL(group.Picture)
			}
			container.Objects[1].(*widget.Label).SetText(group.Name)
			container.Objects[3].(*canvas.Text).Text = group.ID
		},
	)

	gw.widget.OnSelected = func(i widget.ListItemID) {
		state.selected = state.groups[i]
		gw.widget.RefreshItem(i)
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
