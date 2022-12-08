package src

import (
	"fmt"

	"github.com/delucks/go-subsonic"
	"github.com/rivo/tview"
)

func (a *app) buildHeader() tview.Primitive {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)

	a.headerSections = tview.NewTextView().
		SetRegions(true).
		SetChangedFunc(func() {
			a.tv.Draw()
		}).
		SetHighlightedFunc(func(added, _, _ []string) {
			hl := added[0]
			cur, _ := a.pages.GetFrontPage()

			if hl != cur {
				a.switchToPage(hl)
			}
		})
	fmt.Fprintf(a.headerSections, `["artists"]F1: Artists[""] | ["playlists"]F2: Playlists[""] | ["config"]F3: Configuration[""]`)

	a.headerNowPlaying = tview.NewTextView().SetTextAlign(tview.AlignRight)

	flex.AddItem(a.headerSections, 0, 1, false)
	flex.AddItem(a.headerNowPlaying, 0, 1, false)

	a.playQueue.SetOnChangeCallback(func(song *subsonic.Child, isPaused bool) {
		if song != nil {
			symbol := ">"
			if isPaused {
				symbol = "||"
			}
			a.headerNowPlaying.SetText(fmt.Sprintf("%s %s - %s", symbol, song.Title, song.Artist))
		} else {
			a.headerNowPlaying.SetText("Not playing")
		}
	})

	return flex
}
