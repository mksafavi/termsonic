package src

import (
	"fmt"
	"os"

	"github.com/delucks/go-subsonic"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type app struct {
	// General GUI
	tv     *tview.Application
	pages  *tview.Pages
	header *tview.TextView
	footer *tview.TextView
	cfg    *Config

	// Artists panel
	artistsTree *tview.TreeView
	songsList   *tview.List

	// Subsonic variables
	sub *subsonic.Client
}

func Run(cfg *Config) {
	a := &app{
		cfg: cfg,
	}

	a.tv = tview.NewApplication()
	a.pages = tview.NewPages()
	a.footer = tview.NewTextView()

	a.header = tview.NewTextView().
		SetRegions(true).
		SetChangedFunc(func() {
			a.tv.Draw()
		}).
		SetHighlightedFunc(func(added, _, _ []string) {
			hl := added[0]
			cur, _ := a.pages.GetFrontPage()

			if hl != cur {
				switchToPage(a, hl)
			}
		})
	fmt.Fprintf(a.header, `["artists"]F1: Artists[""] | ["playlists"]F2: Playlists[""] | ["config"]F3: Configuration[""]`)

	a.pages.SetBorder(true)
	a.pages.AddPage("config", configPage(a), true, false)
	a.pages.AddPage("artists", artistsPage(a), true, false)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.header, 1, 1, false).
		AddItem(a.pages, 0, 3, true).
		AddItem(a.footer, 1, 1, false)

	if testConfig(a.cfg) != nil {
		switchToPage(a, "config")
	} else {
		a.sub, _ = buildSubsonicClient(a.cfg)
		err := refreshArtists(a)
		if err != nil {
			alert(a, "Could not refresh artists: %v", err)
		} else {
			switchToPage(a, "artists")
		}
	}

	// Keyboard shortcuts
	a.tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			switchToPage(a, "artists")
			return nil
		case tcell.KeyF2:
			switchToPage(a, "playlists")
			return nil
		case tcell.KeyF3:
			switchToPage(a, "config")
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

func switchToPage(a *app, name string) {
	switch name {
	case "artists":
		a.pages.SwitchToPage("artists")
		a.header.Highlight("artists")
	case "playlists":
		a.pages.SwitchToPage("playlists")
		a.header.Highlight("playlists")
	case "config":
		a.pages.SwitchToPage("config")
		a.header.Highlight("config")
	}
}
