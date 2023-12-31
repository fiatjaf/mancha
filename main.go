package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nbd-wtf/go-nostr"
	"golang.org/x/net/context"
)

const (
	APP_TITLE = "Mancha"
	APPID     = "com.nostr.mancha"
)

var baseSize = fyne.Size{Width: 900, Height: 640}

var (
	a    fyne.App
	w    fyne.Window
	k    Keystore
	t    *CustomTheme
	pool = nostr.NewSimplePool(context.Background())
)

// main widgets
var (
	actionsWidget  = makeActionsWidget()
	bottomWidget   = makeBottomWidget()
	messagesWidget = makeMessagesWidget()
	groupsWidget   = makeGroupsWidget()

	emptyRelayListOverlay = container.NewCenter(widget.NewButtonWithIcon("Join Group", theme.StorageIcon(), func() {
		fmt.Println("adding")
		// addRelayDialog(relaysListWidget, chatMessagesListWidget)
	}))
)

func main() {
	a = app.NewWithID(APPID)
	w = a.NewWindow(APP_TITLE)
	t = NewCustomTheme()
	a.Settings().SetTheme(t)
	w.Resize(baseSize)

	// keystore might be using the native keyring or falling back to just a file with a key
	k = startKeystore()

	leftBorderContainer := container.NewBorder(
		nil,
		container.NewPadded(actionsWidget.widget),
		nil,
		nil,
		container.NewStack(
			container.NewPadded(groupsWidget.widget),
			// emptyRelayListOverlay,
		),
	)

	rightBorderContainer := container.NewBorder(
		nil,
		container.NewPadded(bottomWidget),
		nil,
		nil,
		container.NewPadded(messagesWidget.widget),
	)

	splitContainer := container.NewHSplit(leftBorderContainer, rightBorderContainer)
	splitContainer.Offset = 0.35

	w.SetContent(splitContainer)

	go func() {
		for _, group := range loadGroups() {
			state.groups = append(state.groups, &group)
		}
	}()

	w.ShowAndRun()
}
