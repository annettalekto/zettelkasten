package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type viewType struct {
	Text *widget.Entry
	Date *widget.Label
	Name *widget.Entry
}

var view viewType

func refreshTabs() {
	view.Date.SetText(fmt.Sprintf("%v", selectedFile.data.Format("2006-01-02 15:04"))) //d.Format("2006-01-02 15:04")
	view.Name.SetText(selectedFile.title)
	view.Text.SetText(getTextFromFile(selectedFile.filePath))
}

func (v *viewType) viewForm() *fyne.Container { // initViewForm

	// номер карты в шапке?
	v.Date = newFormatLabel(fmt.Sprintf("%v", selectedFile.data))
	d := container.NewBorder(nil, nil, nil, v.Date)
	v.Name = newFormatEntry()
	v.Name.SetText("Имя файла")
	top := container.NewVBox(d, v.Name)

	v.Text = widget.NewMultiLineEntry()
	v.Text.TextStyle.Monospace = true
	v.Text.Wrapping = fyne.TextWrapWord
	v.Text.SetText("Текст")

	bottom := container.NewBorder(nil, nil, nil, widget.NewButton("Редакт.", nil))

	return container.NewBorder(top, bottom, nil, nil, v.Text)
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
