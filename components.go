package main

import (
	"fmt"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type middleClickHandlerEntry struct {
	widget.Entry
}

func (e *middleClickHandlerEntry) MouseUp(ev *desktop.MouseEvent) {
	fmt.Println("click")

	if ev.Button == desktop.MouseButtonTertiary {
		paste := getLinuxPrimaryClipboard()
		c := e.Entry.CursorColumn
		current := e.Entry.Text
		newText := current[0:c] + paste + current[c:]
		e.Entry.Text = newText
		e.Entry.CursorColumn += len(paste)
		e.Entry.SetText(newText)
	}

	if ev.Button == desktop.MouseButtonPrimary {
		setLinuxPrimaryClipboard(e.Entry.SelectedText())
	}

	e.Entry.MouseUp(ev)
}
