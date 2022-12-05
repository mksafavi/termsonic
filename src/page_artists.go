package src

import "github.com/rivo/tview"

func artistsPage(a *app) tview.Primitive {
	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(30, 0).
		SetBorders(true)

	grid.AddItem(tview.NewTextView().SetText("Artist & Album list"), 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(tview.NewTextView().SetText("Song list!"), 0, 1, 1, 2, 0, 0, false)

	return grid
}
