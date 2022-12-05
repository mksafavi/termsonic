package src

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type selection struct {
	entryType string
	id        string
}

func artistsPage(a *app) tview.Primitive {
	grid := tview.NewGrid().
		SetColumns(40, 0).
		SetBorders(true)

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

			loadAlbumInPanel(a, sel.id)
			a.tv.SetFocus(a.songsList)
		})

	a.songsList = tview.NewList()
	a.songsList.ShowSecondaryText(false)

	grid.AddItem(a.artistsTree, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(a.songsList, 0, 1, 1, 2, 0, 0, false)

	return grid
}

func refreshArtists(a *app) error {
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

func loadAlbumInPanel(a *app, id string) error {
	album, err := a.sub.GetMusicDirectory(id)
	if err != nil {
		return err
	}

	a.songsList.SetTitle(album.Name)
	a.songsList.Clear()
	for _, song := range album.Child {
		dur := time.Duration(song.Duration) * time.Second
		a.songsList.AddItem(fmt.Sprintf("%-10s %d - %s", fmt.Sprintf("[%s]", dur.String()), song.Track, song.Title), "", 0, nil)
	}

	return nil
}
