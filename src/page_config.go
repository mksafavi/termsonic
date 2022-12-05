package src

import (
	"fmt"
	"net/http"

	"github.com/delucks/go-subsonic"
	"github.com/rivo/tview"
)

func configPage(a *app) *tview.Form {
	var err error

	form := tview.NewForm().
		AddInputField("Server URL", a.cfg.BaseURL, 40, nil, func(txt string) { a.cfg.BaseURL = txt }).
		AddInputField("Username", a.cfg.Username, 20, nil, func(txt string) { a.cfg.Username = txt }).
		AddPasswordField("Password", a.cfg.Password, 20, '*', func(txt string) { a.cfg.Password = txt }).
		AddButton("Test", func() {
			if err = testConfig(a.cfg); err != nil {
				alert(a, "Could not auth: %v", err)
			} else {
				alert(a, "Success.")
			}
		}).
		AddButton("Save", func() {
			err := a.cfg.Save()
			if err != nil {
				alert(a, "Error saving: %v", err)
				return
			}

			a.sub, err = buildSubsonicClient(a.cfg)
			if err != nil {
				alert(a, "Could not auth: %v", err)
			} else {
				alert(a, "All good!")
			}
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