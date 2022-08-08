package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Zettelkasten")
	w.Resize(fyne.NewSize(800, 600))

	w.SetContent(widget.NewLabel("Let's start..."))
	w.ShowAndRun()
}
