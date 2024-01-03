package main

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ActionsWidget struct {
	widget *widget.Toolbar
}

var _actionsWidget *ActionsWidget

func getActionWidget() *ActionsWidget {
	if _actionsWidget == nil {
		_actionsWidget = makeActionsWidget()
	}
	return _actionsWidget
}

func makeActionsWidget() *ActionsWidget {
	aw := &ActionsWidget{}

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
		widget.NewToolbarAction(theme.ContentAddIcon(), showAddGroupDialog),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			dialog.NewConfirm("Reset local data?", "This will remove all relays and your private key.", func(b bool) {
				if b {
					state.groups = nil
					saveGroups()
					getGroupsWidget().widget.Refresh()
					getMessagesWidget().widget.Refresh()
					k.Erase()
				}
			}, w).Show()
		}),
	)

	return aw
}
