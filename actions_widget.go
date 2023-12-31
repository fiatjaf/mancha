package main

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ActionWidget struct {
	widget *widget.Toolbar
}

func makeActionsWidget() *ActionWidget {
	aw := &ActionWidget{}

	aw.widget = widget.NewToolbar(
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			entry := widget.NewEntry()
			entry.SetPlaceHolder("nsec1...")
			dialog.ShowForm("Import a Nostr Private Key                                             ", "Import", "Cancel", []*widget.FormItem{ // Empty space Hack to make dialog bigger
				widget.NewFormItem("Private Key", entry),
			}, func(b bool) {
				if entry.Text != "" && b {
					err := saveKey(entry.Text) // TODO: Handle Error
					if err != nil {
						fmt.Println("Err saving key: ", err)
					}
				}
			}, w)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.StorageIcon(), func() {
			fmt.Println("toolbar action")
			// addRelayDialog(relaysListWidget, chatMessagesListWidget)
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			dialog.NewConfirm("Reset local data?", "This will remove all relays and your private key.", func(b bool) {
				if b {
					state.groups = nil
					saveGroups(state.groups)
					groupsWidget.widget.Refresh()

					messagesWidget.widget.Refresh()
					k.Erase()
				}
			}, w).Show()
		}),
	)

	return aw
}
