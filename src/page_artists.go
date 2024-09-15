package src

import (
	"fmt"

	"github.com/delucks/go-subsonic"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type selection struct {
	entryType string
	id        string
}

func (a *app) nowPlaying() {
	a.nowPlayingFlex = tview.NewFlex().SetDirection(tview.FlexColumn)
	songText := tview.NewTextView().SetTextAlign(tview.AlignLeft)
	coverArt := tview.NewImage()

	a.nowPlayingFlex.SetBorder(true)
	a.nowPlayingFlex.SetTitleAlign(tview.AlignLeft)
	a.nowPlayingFlex.SetTitle(fmt.Sprintf("paused"))

	a.nowPlayingFlex.AddItem(coverArt, 0, 1, false)
	a.nowPlayingFlex.AddItem(songText, 0, 3, false)

	a.playQueue.SetOnChangeCallback(func(song *subsonic.Child, isPaused bool) {
		if song != nil {
			if isPaused {
				a.nowPlayingFlex.SetTitle(fmt.Sprintf("Paused"))
			} else {
				a.nowPlayingFlex.SetTitle(fmt.Sprintf("Playing"))
			}

			songText.SetText(fmt.Sprintf("Title: %s\nAlbum: %s\nArtist: %s", song.Title, song.Album, song.Artist))
			coverArtImage, err := a.sub.GetCoverArt(song.ID, nil)
			if err == nil {
				coverArt.SetImage(coverArtImage)
				coverArt.SetColors(tview.TrueColor)
				coverArt.SetDithering(tview.DitheringNone)
			}
		} else {
			songText.SetText("")
		}

		a.updatePageQueue()

		// Fix "Now Playing" not always updating
		go a.tv.Draw()
	})
}

func (a *app) artistsPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	left_flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Artist & album list
	root := tview.NewTreeNode("Subsonic server").SetColor(tcell.ColorYellow)
	root.SetSelectable(false)
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
				if node.GetChildren() == nil || len(node.GetChildren()) == 0 {
					artist, err := a.sub.GetArtist(sel.id)
					if err != nil {
						LogErrorf("loading artist '%s': %v", sel.id, err)
						a.alert("Error: %v", err)
						return
					}

					for _, album := range artist.Album {
						subnode := tview.NewTreeNode(album.Name)
						subnode.SetReference(selection{"album", album.ID})
						subnode.SetColor(tcell.ColorBlue)
						subnode.SetSelectable(true)

						node.AddChild(subnode)
					}

				}

				node.SetExpanded(!node.IsExpanded())
				return
			}

			a.loadAlbumInPanel(sel.id)
			a.tv.SetFocus(a.songsList)
		})
	a.artistsTree.SetBorderAttributes(tcell.AttrDim).SetBorder(true)
	a.artistsTree.SetFocusFunc(func() { a.updateFooter() })

	// Songs list for the selected album
	a.songsList = tview.NewList()
	a.songsList.ShowSecondaryText(false).SetHighlightFullLine(true)
	a.songsList.SetBorderAttributes(tcell.AttrDim).SetBorder(true)
	a.songsList.SetFocusFunc(func() { a.updateFooter() })

	// Change the left-right keys to switch between the panels
	a.artistsTree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			a.tv.SetFocus(a.songsList)
			return nil
		}
		return event
	})

	a.songsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			a.tv.SetFocus(a.artistsTree)
			return nil
		}
		return event
	})

	a.setupKeybindings(flex.Box)

	flex.AddItem(a.artistsTree, 0, 1, true)
	flex.AddItem(left_flex, 0, 1, false)
	left_flex.AddItem(a.songsList, 0, 3, false)
	left_flex.AddItem(a.nowPlayingFlex, 0, 1, false)
	return flex
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

			a.artistsTree.GetRoot().AddChild(node)
		}
	}

	a.artistsTree.GetRoot().SetExpanded(true)

	return nil
}

func (a *app) loadAlbumInPanel(id string) error {
	album, err := a.sub.GetAlbum(id)
	if err != nil {
		panic(err)
	}

	a.songsList.Clear()
	a.currentSongs = album.Song
	for _, song := range album.Song {
		txt := fmt.Sprintf("[%02d:%02d:%02d] %d.%d - %s", (song.Duration / 3600), (song.Duration / 60), (song.Duration % 60), song.DiscNumber, song.Track, song.Title)

		a.songsList.AddItem(txt, "", 0, func() {
			sel := a.songsList.GetCurrentItem()
			a.playQueue.Clear()
			for _, s := range a.currentSongs[sel:] {
				a.playQueue.Append(s)
			}
			err := a.playQueue.Play()
			if err != nil {
				a.alert("Error: %v", err)
				LogErrorf("starting playback of album '%s': %v", album.Name, err)
			}
		})
	}

	a.songsList.SetCurrentItem(0)

	return nil
}
