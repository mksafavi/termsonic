package src

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type selection struct {
	entryType string
	id        string
}

func (a *app) artistsPage() tview.Primitive {
	grid := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Artist & album list
	root := tview.NewTreeNode("Subsonic server").SetColor(tcell.ColorYellow)
	a.artistsTree = tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root).
		SetPrefixes([]string{"", " ", " "}).
		SetSelectedFunc(func(node *tview.TreeNode) {
			if node.GetReference() == nil {
				return
			}

			sel := node.GetReference().(selection)
			if sel.entryType == "artist" {
				node.SetExpanded(!node.IsExpanded())
				return
			}

			a.loadAlbumInPanel(sel.id)
			a.tv.SetFocus(a.songsList)
			a.updateFooter()
		})
	a.artistsTree.SetBorderAttributes(tcell.AttrDim).SetBorder(true)

	// Songs list for the selected album
	a.songsList = tview.NewList()
	a.songsList.ShowSecondaryText(false).SetHighlightFullLine(true)
	a.songsList.SetBorderAttributes(tcell.AttrDim).SetBorder(true)

	// Change the left-right keys to switch between the panels
	a.artistsTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			a.tv.SetFocus(a.songsList)
			a.updateFooter()
			return nil
		}
		return event
	})

	a.songsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			a.tv.SetFocus(a.artistsTree)
			a.updateFooter()
			return nil
		}
		return event
	})

	a.setupKeybindings(grid.Box)

	grid.AddItem(a.artistsTree, 0, 1, true)
	grid.AddItem(a.songsList, 0, 1, false)

	return grid
}

func (a *app) refreshArtists() error {
	artistsID3, err := a.sub.GetArtists(nil)
	if err != nil {
		return err
	}

	a.artistsTree.GetRoot().ClearChildren()
	for _, index := range artistsID3.Index {
		for _, artist := range index.Artist {
			node := tview.NewTreeNode(artist.Name)
			node.SetReference(selection{"artist", artist.ID})
			node.SetColor(tcell.ColorRed)
			node.SetSelectable(true)
			node.SetExpanded(false)

			albums, err := a.sub.GetMusicDirectory(artist.ID)
			if err != nil {
				return err
			}

			for _, album := range albums.Child {
				subnode := tview.NewTreeNode(album.Title)
				subnode.SetReference(selection{"album", album.ID})
				subnode.SetColor(tcell.ColorBlue)
				subnode.SetSelectable(true)

				node.AddChild(subnode)
			}

			a.artistsTree.GetRoot().AddChild(node)
		}
	}

	a.artistsTree.GetRoot().SetExpanded(true)

	return nil
}

func (a *app) loadAlbumInPanel(id string) error {
	album, err := a.sub.GetMusicDirectory(id)
	if err != nil {
		return err
	}

	a.songsList.Clear()
	a.currentSongs = album.Child
	for _, song := range album.Child {
		txt := fmt.Sprintf("%-2d - %s", song.Track, song.Title)

		a.songsList.AddItem(txt, "", 0, func() {
			sel := a.songsList.GetCurrentItem()
			a.playQueue.Clear()
			for _, s := range a.currentSongs[sel:] {
				a.playQueue.Append(s)
			}
			err := a.playQueue.Play()
			if err != nil {
				a.alert("Error: %v", err)
			}
		})
	}

	a.songsList.SetCurrentItem(0)

	return nil
}
