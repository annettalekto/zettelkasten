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
	name.SetText("Имя файла")
	top := container.NewVBox(d, name)

	text := widget.NewMultiLineEntry()
	text.SetText("Текст")

	bottom := container.NewBorder(nil, nil, nil, widget.NewButton("Редакт.", nil))

	return container.NewBorder(top, bottom, nil, nil, text)
}

func addInfoForm() *fyne.Container {

	name := newFormatEntry()
	name.SetText("Имя файла")

	tegs := widget.NewMultiLineEntry()
	tegs.SetText("#тег1\n#тег2")
	binds := widget.NewMultiLineEntry()
	binds.SetText("Связное")
	source := widget.NewMultiLineEntry()
	source.SetText("Источники")
	box := container.NewGridWithColumns(1, tegs, binds, source)

	return container.NewBorder(name, nil, nil, nil, box)
}

func sourceInfoForm() *fyne.Container {
	source := widget.NewMultiLineEntry()
	source.SetText("Источник")

	quotation := widget.NewMultiLineEntry()
	quotation.SetText("Цитата")

	return container.NewGridWithColumns(1, source, quotation) // todo: разделить
}

func commentForm() *fyne.Container {
	comment := widget.NewMultiLineEntry()
	comment.SetText("Комментарий")

	return container.NewGridWithColumns(1, comment)
}
