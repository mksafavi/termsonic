package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/delucks/go-subsonic"
	"github.com/rivo/tview"
)

var (
	url      = ""
	username = ""
	password = ""

	sub *subsonic.Client = nil

	// GUI
	app   *tview.Application
	pages *tview.Pages
)

func main() {
	loadConfig()

	app = tview.NewApplication()

	pages = tview.NewPages()

	showConfig := sub == nil
	pages.AddPage("page-config", configPage(), true, showConfig)

	if !showConfig {
		pages.AddPage("page-main", mainView(), true, true)
	}

	if err := app.SetRoot(pages, true).EnableMouse(true).SetFocus(pages).Run(); err != nil {
		fmt.Printf("Error running TermSonic: %v", err)
		os.Exit(1)
	}
}

func configPage() *tview.Form {
	form := tview.NewForm().
		AddInputField("Server URL", url, 40, nil, func(txt string) { url = txt }).
		AddInputField("Username", username, 20, nil, func(txt string) { username = txt }).
		AddPasswordField("Password", password, 20, '*', func(txt string) { password = txt }).
		AddButton("Test", func() {
			tmpSub := &subsonic.Client{
				Client:       http.DefaultClient,
				BaseUrl:      url,
				User:         username,
				ClientName:   "termsonic",
				PasswordAuth: true,
			}

			if err := tmpSub.Authenticate(password); err != nil {
				alert("Could not auth: %v", err)
			} else {
				sub = tmpSub

				alert("Success.")
			}
		}).
		AddButton("Save", nil)
	return form
}

func mainView() *tview.Grid {
	grid := tview.NewGrid().
		SetRows(2, 0).
		SetColumns(30, 0).
		SetBorders(true)

	grid.AddItem(tview.NewTextView().SetText("Artist & Album list"), 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(tview.NewTextView().SetText("Song list!"), 1, 1, 1, 2, 0, 0, false)

	return grid
}

func alert(format string, params ...interface{}) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf(format, params...)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			pages.RemovePage("alert")
		})

	pages.AddPage("alert", modal, true, true)
}

func loadConfig() {
	url = "http://music.nuc.local"
	username = "admin"
	password = "admin"

	sub = &subsonic.Client{
		Client:       http.DefaultClient,
		User:         username,
		BaseUrl:      url,
		ClientName:   "termsonic",
		PasswordAuth: true,
	}

	if err := sub.Authenticate(password); err != nil {
		sub = nil
	}
}

func saveConfig() {
}
