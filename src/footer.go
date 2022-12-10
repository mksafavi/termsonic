package src

func (a *app) updateFooter() {
	switch a.headerSections.GetHighlights()[0] {
	case "artists":
		if a.tv.GetFocus() == a.artistsTree {
			a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause   [blue]e:[yellow] Play album last   [blue]n:[yellow] Play album next")
		} else if a.tv.GetFocus() == a.songsList {
			a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause   [blue]e:[yellow] Play song last   [blue]n:[yellow] Play song next")
		}
	case "playqueue":
		a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause   [blue]d:[yellow] Remove   [blue]j:[yellow] Move up   [blue]k:[yellow] Move down")
	case "playlists":
		a.footer.SetText("")
	case "config":
		a.footer.SetText("")
	}
}
