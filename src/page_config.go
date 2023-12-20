package src

import (
	"fmt"
	"net/http"

	"github.com/delucks/go-subsonic"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) configPage() *tview.Form {
	var err error

	form := tview.NewForm().
		AddInputField("Server URL", a.cfg.BaseURL, 40, nil, func(txt string) { a.cfg.BaseURL = txt }).
		AddInputField("Username", a.cfg.Username, 20, nil, func(txt string) { a.cfg.Username = txt }).
		AddPasswordField("Password", a.cfg.Password, 20, '*', func(txt string) { a.cfg.Password = txt }).
		AddButton("Test", func() {
			if err = testConfig(a.cfg); err != nil {
				a.alert("Could not auth: %v", err)
			} else {
				a.alert("Success.")
			}
		}).
		AddButton("Save", func() {
			err := a.cfg.Save()
			if err != nil {
				a.alert("Error saving: %v", err)
				return
			}

			a.sub, err = buildSubsonicClient(a.cfg)
			if err != nil {
				a.alert("Could not auth: %v", err)
			} else {
				a.playQueue.SetClient(a.sub)
				a.alert("All good!")
			}
		})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlR {
			if a.sub == nil {
				return nil
			}

			if err := a.refreshArtists(); err != nil {
				a.alert("Error: %v", err)
				LogErrorf("refreshing artists following Ctrl+R: %v", err)
				return nil
			}

			if err := a.refreshPlaylists(); err != nil {
				a.alert("Error: %v", err)
				LogErrorf("refreshing playlists following Ctrl+R: %v", err)
				return nil
			}

			a.alert("Refreshed successfully")

			return nil
		}

		return event
	})

	return form
}

func testConfig(cfg *Config) error {
	if cfg.BaseURL == "" {
		return fmt.Errorf("empty base URL")
	}

	if cfg.Username == "" {
		return fmt.Errorf("empty username")
	}

	if cfg.Password == "" {
		return fmt.Errorf("empty password")
	}

	_, err := buildSubsonicClient(cfg)
	return err
}

func buildSubsonicClient(cfg *Config) (*subsonic.Client, error) {
	tmpSub := &subsonic.Client{
		Client:       http.DefaultClient,
		BaseUrl:      cfg.BaseURL,
		User:         cfg.Username,
		ClientName:   "termsonic",
		PasswordAuth: true,
	}

	err := tmpSub.Authenticate(cfg.Password)
	if err != nil {
		return nil, err
	}

	return tmpSub, nil
}
