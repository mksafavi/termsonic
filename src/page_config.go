package src

import (
	"net/http"

	"github.com/delucks/go-subsonic"
	"github.com/rivo/tview"
)

func configPage(a *app) *tview.Form {
	form := tview.NewForm().
		AddInputField("Server URL", a.cfg.BaseURL, 40, nil, func(txt string) { a.cfg.BaseURL = txt }).
		AddInputField("Username", a.cfg.Username, 20, nil, func(txt string) { a.cfg.Username = txt }).
		AddPasswordField("Password", a.cfg.Password, 20, '*', func(txt string) { a.cfg.Password = txt }).
		AddButton("Test", func() {
			tmpSub := &subsonic.Client{
				Client:       http.DefaultClient,
				BaseUrl:      a.cfg.BaseURL,
				User:         a.cfg.Username,
				ClientName:   "termsonic",
				PasswordAuth: true,
			}

			if err := tmpSub.Authenticate(a.cfg.Password); err != nil {
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

			a.sub = &subsonic.Client{
				Client:       http.DefaultClient,
				BaseUrl:      a.cfg.BaseURL,
				User:         a.cfg.Username,
				ClientName:   "termsonic",
				PasswordAuth: true,
			}
			if err := a.sub.Authenticate(a.cfg.Password); err != nil {
				alert(a, "Could not auth: %v", err)
			} else {
				alert(a, "All good!")
			}
		})
	return form
}
