package main

import (
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func newEnhancedEntry() *EnhancedEntry {
	ee := &EnhancedEntry{}
	ee.ExtendBaseWidget(ee)
	return ee
}

type EnhancedEntry struct {
	widget.Entry
}

func (e *EnhancedEntry) MouseUp(ev *desktop.MouseEvent) {
	if ev.Button == desktop.MouseButtonTertiary {
		paste := getLinuxPrimaryClipboard()
		c := e.Entry.CursorColumn
		if r := e.Entry.CursorRow; r > 0 {
			c = e.Entry.TextPosFromRowCol(r, c)
		}

		current := e.Entry.Text
		newText := current[0:c] + paste + current[c:]
		e.Entry.Text = newText
		e.Entry.CursorColumn += len(paste)
		e.Entry.SetText(newText)
	}

	e.Entry.MouseUp(ev)
}
