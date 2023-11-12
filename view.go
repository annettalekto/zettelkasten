package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func viewForm() *fyne.Container {

	// номер карты в шапке?
	date := newFormatLabel("дата: 00.00.0000")
	d := container.NewBorder(nil, nil, nil, date)
	name := newFormatEntry()
	top := container.NewVBox(d, name)

	text := widget.NewMultiLineEntry()

	bottom := container.NewBorder(nil, nil, nil, widget.NewButton("Редакт.", nil))

	return container.NewBorder(top, bottom, nil, nil, text)
}
