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
	a.footer = tview.NewTextView().
		SetDynamicColors(true)

	a.header = tview.NewTextView().
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
	fmt.Fprintf(a.header, `["artists"]F1: Artists[""] | ["playlists"]F2: Playlists[""] | ["config"]F3: Configuration[""]`)

	a.pages.SetBorder(true)
	a.pages.AddPage("config", a.configPage(), true, false)
	a.pages.AddPage("artists", a.artistsPage(), true, false)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.header, 1, 1, false).
		AddItem(a.pages, 0, 3, true).
		AddItem(a.footer, 1, 1, false)

	if testConfig(a.cfg) != nil {
		a.switchToPage("config")
	} else {
		a.sub, _ = buildSubsonicClient(a.cfg)
		err := a.refreshArtists()
		if err != nil {
			a.alert("Could not refresh artists: %v", err)
		} else {
			a.switchToPage("artists")
		}
	}

	// Keyboard shortcuts
	a.tv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			a.switchToPage("artists")
			return nil
		case tcell.KeyF2:
			a.switchToPage("playlists")
			return nil
		case tcell.KeyF3:
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
		a.pages.SwitchToPage("artists")
		a.header.Highlight("artists")
		a.tv.SetFocus(a.artistsTree)
	case "playlists":
		a.pages.SwitchToPage("playlists")
		a.header.Highlight("playlists")
	case "config":
		a.pages.SwitchToPage("config")
		a.header.Highlight("config")
	}

	a.updateFooter()
}
