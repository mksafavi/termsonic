package src

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) playlistsPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)

	a.playlistsList = tview.NewList().
		SetMainTextColor(tcell.ColorRed).
		SetHighlightFullLine(true).
		ShowSecondaryText(false)
	a.playlistsList.SetBorder(true).SetBorderAttributes(tcell.AttrDim)

	a.playlistSongs = tview.NewList().
		SetHighlightFullLine(true).
		ShowSecondaryText(false)
	a.playlistSongs.SetBorder(true).SetBorderAttributes(tcell.AttrDim)

	// Change the left-right keys to switch between the panels
	a.playlistsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			a.tv.SetFocus(a.playlistSongs)
			a.updateFooter()
			return nil
		}
		return event
	})
	a.playlistSongs.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			a.tv.SetFocus(a.playlistsList)
			a.updateFooter()
			return nil
		}
		return event
	})

	// Setup e & n keybinds
	a.setupMusicControlKeys(flex.Box)

	flex.AddItem(a.playlistsList, 0, 1, false)
	flex.AddItem(a.playlistSongs, 0, 1, false)

	return flex
}

func (a *app) refreshPlaylists() error {
	playlists, err := a.sub.GetPlaylists(nil)
	if err != nil {
		return err
	}

	a.allPlaylists = playlists

	a.playlistsList.Clear()
	for _, pl := range playlists {
		id := pl.ID
		a.playlistsList.AddItem(pl.Name, "", 0, func() {
			a.loadPlaylist(id)
			a.tv.SetFocus(a.playlistSongs)
			a.updateFooter()
		})
	}

	a.playlistsList.SetCurrentItem(0)

	return nil
}

func (a *app) loadPlaylist(id string) error {
	a.playlistSongs.Clear()
	pl, err := a.sub.GetPlaylist(id)
	if err != nil {
		return err
	}

	a.currentPlaylist = pl

	for _, s := range a.currentPlaylist.Entry {
		a.playlistSongs.AddItem(fmt.Sprintf("%s - %s", s.Title, s.Artist), "", 0, func() {
			sel := a.playlistSongs.GetCurrentItem()
			a.playQueue.Clear()
			for _, s := range a.currentPlaylist.Entry[sel:] {
				a.playQueue.Append(s)
			}

			if err := a.playQueue.Play(); err != nil {
				a.alert("Error: %v", err)
			}
		})
	}

	return nil
}
