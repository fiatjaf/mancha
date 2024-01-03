package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
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

	rapidClicks int
}

func (e *EnhancedEntry) MouseDown(ev *desktop.MouseEvent) {
	if ev.Button == desktop.MouseButtonTertiary && !e.Entry.Disabled() {
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

	if ev.Button == desktop.MouseButtonPrimary {
		e.rapidClicks++
		go func() {
			time.Sleep(time.Millisecond * 350)
			e.rapidClicks--
		}()
		if e.rapidClicks == 3 {
			fmt.Println("3 rapid clicks, must select all")
		}
	}

	e.Entry.MouseDown(ev)
}

func (e *EnhancedEntry) KeyDown(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn {
		e.Entry.OnSubmitted(e.Entry.Text)
	}

	e.Entry.KeyDown(ev)
}
