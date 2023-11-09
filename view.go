package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func viewForm() *fyne.Container {

	// название сверху и дата
	// текст
	// кнопка редактирования
	box := container.NewHBox(widget.NewLabel("lasdk"))
	return box
}
