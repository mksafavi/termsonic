package src

func (a *app) updateFooter() {
	switch a.header.GetHighlights()[0] {
	case "artists":
		switch a.tv.GetFocus() {
		case a.artistsTree:
			a.footer.SetText("Artists: [blue]Up/Down:[yellow] Move selection    [blue]Space:[yellow] Select entry")
		case a.songsList:
			a.footer.SetText("Songs:   [blue]Up/Down:[yellow] Move selection    [blue]Space:[yellow] Play")
		}
	case "playlists":
		a.footer.SetText("Come back later!")
	case "config":
		a.footer.SetText("Configuration page")
	}
}
