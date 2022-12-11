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
		a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause   [blue]d:[yellow] Remove   [blue]j:[yellow] Move up   [blue]k:[yellow] Move down  [blue]r:[yellow] Shuffle")
	case "playlists":
		if a.tv.GetFocus() == a.playlistsList {
			a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause   [blue]e:[yellow] Play playlist last   [blue]n:[yellow] Play playlist next   [blue]r:[yellow] Random")
		} else if a.tv.GetFocus() == a.playlistSongs {
			a.footer.SetText("[blue]l:[yellow] Next song   [blue]p:[yellow] Toggle pause   [blue]e:[yellow] Play song last   [blue]n:[yellow] Play song next")
		}
	case "config":
		a.footer.SetText("[blue]R:[yellow] Refresh")
	}
}
