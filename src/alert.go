package src

import (
	"fmt"

	"github.com/rivo/tview"
)

func alert(a *app, format string, params ...interface{}) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf(format, params...)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			a.pages.RemovePage("alert")
		})

	a.pages.AddPage("alert", modal, true, true)
}
