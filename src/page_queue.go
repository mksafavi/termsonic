package src

import (
	"fmt"

	"github.com/rivo/tview"
)

func (a *app) queuePage() tview.Primitive {
	a.playQueueList = tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	a.setupKeybindings(a.playQueueList.Box)

	a.updatePageQueue()

	return a.playQueueList
}

func (a *app) updatePageQueue() {
	sel := a.playQueueList.GetCurrentItem()
	a.playQueueList.Clear()

	for _, song := range a.playQueue.GetSongs() {
		ownSong := *song
		a.playQueueList.AddItem(fmt.Sprintf("* %s - %s", song.Title, song.Artist), "", 0, func() {
			a.playQueue.SkipTo(&ownSong)
			a.playQueueList.SetCurrentItem(0)
		})
	}

	if sel < a.playQueueList.GetItemCount() {
		a.playQueueList.SetCurrentItem(sel)
	} else {
		a.playQueueList.SetCurrentItem(a.playQueueList.GetItemCount() - 1)
	}
}
