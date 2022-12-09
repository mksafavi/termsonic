package src

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) setupMusicControlKeys(p *tview.Box) {
	// Add 'k' and 'l' key bindings
	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'l' {
			a.playQueue.Next()
			return nil
		}

		if event.Rune() == 'k' {
			a.playQueue.TogglePause()
			return nil
		}
		return event
	})
}
