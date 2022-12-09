package src

func (a *app) updateFooter() {
	switch a.headerSections.GetHighlights()[0] {
	case "artists":
		a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause")
	case "playqueue":
		a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause   [blue]d:[yellow] Remove   [blue]j:[yellow] Move up   [blue]k:[yellow] Move down")
	case "playlists":
		a.footer.SetText("")
	case "config":
		a.footer.SetText("")
	}
}
