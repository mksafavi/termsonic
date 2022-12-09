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

		if event.Rune() == 'p' {
			a.playQueue.TogglePause()
			return nil
		}

		if a.tv.GetFocus() == a.playQueueList {
			if event.Rune() == 'd' {
				sel := a.playQueueList.GetCurrentItem()
				err := a.playQueue.RemoveSong(sel)
				if err != nil {
					a.alert("Error: %v", err)
				}
			} else if event.Rune() == 'k' {
				sel := a.playQueueList.GetCurrentItem()
				if sel == a.playQueueList.GetItemCount()-1 {
					return nil
				}
				err := a.playQueue.Switch(sel, sel+1)
				if err != nil {
					a.alert("Error: %v", err)
				}

				a.playQueueList.SetCurrentItem(sel + 1)

				return nil
			} else if event.Rune() == 'j' {
				sel := a.playQueueList.GetCurrentItem()
				if sel == 0 {
					return nil
				}
				err := a.playQueue.Switch(sel, sel-1)
				if err != nil {
					a.alert("Error: %v", err)
				}

				a.playQueueList.SetCurrentItem(sel - 1)

				return nil
			}
		}

		return event
	})
}
