package src

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"git.sixfoisneuf.fr/termsonic/music"
	"github.com/delucks/go-subsonic"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	// General GUI
	tv               *tview.Application
	pages            *tview.Pages
	headerSections   *tview.TextView
	headerNowPlaying *tview.TextView
	footer           *tview.TextView
	cfg              *Config

	// Artists page
	artistsLoaded bool
	artistsTree   *tview.TreeView
	songsList     *tview.List
	currentSongs  []*subsonic.Child

	// Play queue page
	playQueueList *tview.List

	// Playlist page
	playlistsLoaded bool
	playlistsList   *tview.List
	playlistSongs   *tview.List
	allPlaylists    []*subsonic.Playlist
	currentPlaylist *subsonic.Playlist

	// Subsonic variables
	sub       *subsonic.Client
	playQueue *music.Queue
}

func Run(cfg *Config) {
	a := &app{
		cfg:       cfg,
		playQueue: music.NewQueue(nil),
	}

	a.tv = tview.NewApplication()
	a.pages = tview.NewPages()
	a.footer = tview.NewTextView().
		SetDynamicColors(true)

	a.pages.SetBorder(true)
	a.pages.AddPage("config", a.configPage(), true, false)
	a.pages.AddPage("artists", a.artistsPage(), true, false)
	a.pages.AddPage("playqueue", a.queuePage(), true, false)
	a.pages.AddPage("playlists", a.playlistsPage(), true, false)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.buildHeader(), 1, 1, false).
		AddItem(a.pages, 0, 3, true).
		AddItem(a.footer, 1, 1, false)

	if testConfig(a.cfg) != nil {
		a.switchToPage("config")
	} else {
		a.sub, _ = buildSubsonicClient(a.cfg)
		a.playQueue.SetClient(a.sub)

		fmt.Printf("Loading artists...")
		if err := a.refreshArtists(); err != nil {
			fmt.Println("ERR")
			a.alert("Loading artists: %v", err)
		} else {
			fmt.Println("OK")
			a.artistsLoaded = true
		}

		fmt.Printf("Loading playlists...")
		if err := a.refreshPlaylists(); err != nil {
			fmt.Println("ERR")
			a.alert("Loading playlists: %v", err)
		} else {
			fmt.Println("OK")
			a.playlistsLoaded = true
		}

		a.switchToPage("artists")
	}

	// Keyboard shortcuts
	a.tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			a.switchToPage("artists")
			return nil
		case tcell.KeyF2:
			a.switchToPage("playqueue")
			return nil
		case tcell.KeyF3:
			a.switchToPage("playlists")
			return nil
		case tcell.KeyF4:
			a.switchToPage("config")
			return nil
		}

		switch event.Rune() {
		case 'q':
			a.tv.Stop()
		}

		return event
	})

	if err := a.tv.SetRoot(mainLayout, true).EnableMouse(true).SetFocus(mainLayout).Run(); err != nil {
		fmt.Printf("Error running termsonic: %v", err)
		os.Exit(1)
	}
}

func (a *app) switchToPage(name string) {
	switch name {
	case "artists":
		if !a.artistsLoaded {
			if err := a.refreshArtists(); err != nil {
				a.alert("Error: %v", err)
			}
			a.artistsLoaded = true
		}
		a.pages.SwitchToPage("artists")
		a.headerSections.Highlight("artists")
		a.tv.SetFocus(a.artistsTree)
		a.pages.SetBorder(false)
	case "playqueue":
		a.pages.SwitchToPage("playqueue")
		a.headerSections.Highlight("playqueue")
		a.tv.SetFocus(a.playQueueList)
		a.pages.SetBorder(true)
	case "playlists":
		if !a.playlistsLoaded {
			if err := a.refreshPlaylists(); err != nil {
				a.alert("Error: %v", err)
			}
			a.playlistsLoaded = true
		}
		a.pages.SwitchToPage("playlists")
		a.headerSections.Highlight("playlists")
		a.tv.SetFocus(a.playlistsList)
		a.pages.SetBorder(false)
	case "config":
		a.pages.SwitchToPage("config")
		a.headerSections.Highlight("config")
		a.pages.SetBorder(true)
	}

	a.updateFooter()
}

func randomize(t []*subsonic.Child) []*subsonic.Child {
	t2 := make([]*subsonic.Child, len(t))
	copy(t2, t)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(t2), func(i, j int) { t2[i], t2[j] = t2[j], t2[i] })

	return t2
}
