package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func viewForm() *fyne.Container {

	// название сверху и дата
	name := newFormatEntry()
	date := newFormatLabel("дата:")
	test := newFormatLabelAndEntry("ggg")

	box := container.NewVBox(name, date, test)

	// текст
	// кнопка редактирования

	// box := container.NewHBox(widget.NewLabel("lasdk"))
	return box
}
