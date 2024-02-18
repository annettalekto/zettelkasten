package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type elmFormType struct {
	id        *numericalEntry
	FilePath  *widget.Entry
	Title     *widget.Entry
	Tags      *widget.Entry
	Binds     *widget.Entry
	Source    *widget.Entry
	Text      *widget.Entry
	Date      *widget.Label
	Quotation *widget.Entry
	Comment   *widget.Entry
}

var elmForm elmFormType

func (e *elmFormType) getDataNewCard() (ztc ztcBasicsType) {

	ztc.id = e.id.Text
	// ztc.filePath =  // todo (путь + номер + титл)
	ztc.title = e.Title.Text
	ztc.tags = strings.Split(e.Tags.Text, "\n") // clear
	ztc.bind = strings.Split(e.Binds.Text, "\n")
	// ztc.bindNumbers
	ztc.source = strings.Split(e.Source.Text, "\n")
	// ztc.sourceNumber

	// ztc.data

	// quotation
	// comment

	return ztc
}

func (e *elmFormType) nameCardForm() *fyne.Container {

	e.id = newNumericalEntry() //widget.NewEntry()
	ent1 := container.NewBorder(nil, nil, widget.NewLabel("Номер карт.:"), nil, e.id)
	e.FilePath = widget.NewEntry()
	ent2 := container.NewBorder(nil, nil, widget.NewLabel("Имя файла:    "), nil, e.FilePath)
	// + дополнительно в лейбл вывести полный путь с полным названием (путь + номер + титл)
	okBtn := widget.NewButton("Создать карточку", func() {
		// todo: вызвать сбор всех инфы
		ztc := e.getDataNewCard() // todo: отладить
		fmt.Println(ztc.tags)
	})

	return container.NewVBox(ent1, ent2, okBtn)
}

func (e *elmFormType) textForm() *fyne.Container {

	e.Date = newFormatLabel(fmt.Sprintf("%v", selectedFile.data))
	e.Title = newFormatEntry()
	e.Text = newText()
	e.Text.SetText("<Текст>")

	return container.NewBorder(
		e.Title,
		container.NewBorder(nil, nil, e.Date, nil),
		nil,
		nil,
		e.Text)
}

func (e *elmFormType) addInfoForm() *fyne.Container {

	e.Tags = widget.NewMultiLineEntry()
	e.Binds = widget.NewMultiLineEntry()

	return container.NewBorder(nil, nil, nil, nil,
		container.NewGridWithColumns(1,
			container.NewBorder(widget.NewLabel("Теги:"), nil, nil, nil, e.Tags),
			container.NewBorder(widget.NewLabel("Связное:"), nil, nil, nil, e.Binds)))
}

func (e *elmFormType) sourceForm() *fyne.Container {
	e.Source = newText()
	e.Quotation = newText()

	return container.NewGridWithColumns(1,
		container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, e.Source),
		container.NewBorder(widget.NewLabel("Цитата:"), nil, nil, nil, e.Quotation))
}

func (e *elmFormType) commentForm() *fyne.Container {
	e.Comment = newText()

	return container.NewGridWithColumns(1,
		container.NewBorder(widget.NewLabel("Комментарий:"), nil, nil, nil, e.Comment))
}

func (e *elmFormType) refreshTabs(z ztcBasicsType) {

	e.Date.SetText(fmt.Sprintf("%v", z.data.Format("2006-01-02 15:04")))
	e.Title.SetText(z.title)
	e.Tags.SetText(formatSlice(z.tags))
	e.Binds.SetText(formatSlice(z.bind))
	e.Source.SetText(formatSlice(z.source))
	e.Text.SetText(getTextFromFile(z.filePath))
	e.Quotation.SetText(getQuotationFromFile(z.filePath))
	e.Comment.SetText(getCommentFromFile(z.filePath))
}

/*
func (e *elmFormType) editForm() { // открывать только для текущего файла
	w := fyne.CurrentApp().NewWindow("0")
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()

	// todo: оформить по типу шаблонного файла
	tagBox := container.NewBorder(widget.NewLabel("Теги:"), nil, nil, nil, e.Tags)
	dindsBox := container.NewBorder(widget.NewLabel("Связное:"), nil, nil, nil, e.Binds)
	// номер и название источников отдельно? как удобнее считывать и составлять итоговый файл
	sourceBox := container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, e.Source)
	quotationBox := container.NewBorder(widget.NewLabel("Цитата:"), nil, nil, nil, e.Quotation)
	commentBox := container.NewBorder(widget.NewLabel("Комментарий:"), nil, nil, nil, e.Comment)

	bottom := container.NewGridWithColumns(1, tagBox, dindsBox, sourceBox, quotationBox, commentBox)
	box := container.NewBorder(e.Name, bottom, nil, nil, e.Text)

	w.SetContent(box)
	w.Show()
}
*/
