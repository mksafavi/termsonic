package src

func (a *app) updateFooter() {
	switch a.headerSections.GetHighlights()[0] {
	case "artists":
		a.footer.SetText("[blue]l:[yellow] Next song    [blue]k:[yellow] Toggle pause")
	case "playqueue":
		a.footer.SetText("[blue]l:[yellow] Next song    [blue]k:[yellow] Toggle pause")
	case "playlists":
		a.footer.SetText("Come back later!")
	case "config":
		a.footer.SetText("Configuration page")
	}
}
