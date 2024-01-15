package main

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v3"
	"github.com/nbd-wtf/go-nostr"
	"golang.org/x/net/context"
)

const (
	APP_TITLE = "Mancha"
	APPID     = "com.fiatjaf.mancha"
)

var (
	k    Keystore
	t    *CustomTheme
	pool = nostr.NewSimplePool(context.Background())
)

func main() {
	app := gtk.NewApplication(APPID, gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })
}

func activate(app *gtk.Application) {
	topLabel := gtk.NewLabel("Text set by initializer")
	topLabel.SetVExpand(true)
	topLabel.SetHExpand(true)

	bottomLabel := gtk.NewLabel("Text set by initializer")
	bottomLabel.SetVExpand(true)
	bottomLabel.SetHExpand(true)

	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.Append(topLabel)
	box.Append(bottomLabel)

	window := gtk.NewApplicationWindow(app)
	window.SetTitle(APP_TITLE)
	window.SetChild(box)
	window.SetDefaultSize(900, 640)
	window.Show()

	go func() {
		var ix int
		for t := range time.Tick(time.Second) {
			// Make a copy of the state so we can reference it in the closure.
			currentTime := t
			currentIx := ix

			ix++

			glib.IdleAdd(func() {
				topLabel.SetLabel(fmt.Sprintf("Set a label %d time(s)!", currentIx))
				bottomLabel.SetLabel(fmt.Sprintf(
					"Last updated at %s.",
					currentTime.Format(time.StampMilli),
				))
			})
		}
	}()

	go func() {
		loadPeople()
		loadGroups()
		for _, group := range state.groups {
			go func(group *Group) {
				err := group.startListening()
				if err != nil {
					log.Println(err)
				}
			}(group)

			// getGroupsWidget().overlay.Hide()
		}

		time.Sleep(time.Second * 1)
		// getInputWidget().onMount()
	}()
}
