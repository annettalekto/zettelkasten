package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Zettelkasten")
	w.Resize(fyne.NewSize(700, 500))

	// прога для работы с сиситемой Zettelkasten
	// будет сохранять файлы в определенном виде, ну и читать их
	// тест git
	// конфликт
	// merge конфликт = поправлено

	w.SetContent(widget.NewLabel("Let's start!"))
	w.ShowAndRun()
}
