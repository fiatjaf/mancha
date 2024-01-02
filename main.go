package main

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
		container.NewPadded(getActionWidget().widget),
		nil,
		nil,
		container.NewStack(
			container.NewPadded(getGroupsWidget().widget),
			getGroupsWidget().overlay,
		),
	)

	rightBorderContainer := container.NewBorder(
		nil,
		container.NewPadded(getInputWidget().widget),
		nil,
		nil,
		container.NewPadded(getMessagesWidget().widget),
	)

	splitContainer := container.NewHSplit(leftBorderContainer, rightBorderContainer)
	splitContainer.Offset = 0.35

	w.SetContent(splitContainer)

	go func() {
		loadGroups()
		for _, group := range state.groups {
			go func(group *Group) {
				err := group.startListening()
				if err != nil {
					log.Println(err)
				}
			}(group)
		}

		time.Sleep(time.Second * 1)
		getInputWidget().onMount()
	}()

	w.ShowAndRun()
}
