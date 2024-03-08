package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type elmFormType struct {
	id           *numericalEntry
	FileName     *widget.Entry
	Title        *widget.Entry
	Tags         *widget.Entry
	BindNumbers  *numericalEntry
	Binds        *widget.Entry
	SourceNumber *numericalEntry
	Source       *widget.Entry
	Text         *widget.Entry
	Date         *widget.Entry
	Quotation    *widget.Entry
	Comment      *widget.Entry
}

var view elmFormType

func (e *elmFormType) getDataNewCard() (ztc ztcBasicsType) {

	ztc.id = e.id.Text
	ztc.title = e.Title.Text
	// ztc.filePath =  // todo (путь + номер + титл)
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
	e.id = newNumericalEntry()
	ent1 := container.NewBorder(nil, nil, widget.NewLabel("Номер карт.:"), nil, e.id)
	e.id.Entry.OnChanged = func(s string) {
		e.FileName.SetText(fmt.Sprintf("%s - %s", s, e.Title.Text))
	}
	e.Title = newFormatEntry()
	title := container.NewBorder(nil, nil, widget.NewLabel("Тема:"), nil, e.Title)
	e.Title.OnChanged = func(s string) {
		e.FileName.SetText(fmt.Sprintf("%s - %s", e.id.Text, s))
	}
	e.FileName = widget.NewEntry() // без папки
	ent2 := container.NewBorder(nil, nil, widget.NewLabel("Имя файла:"), nil, e.FileName)
	top := container.NewVBox(ent1, title, ent2)

	e.Text = newText()

	e.Date = newFormatEntry()
	e.Date.SetText(fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04")))
	// + дополнительно в лейбл вывести полный путь с полным названием (путь + номер + титл)
	CreateButton := widget.NewButton("Создать карточку", func() {
		// todo: вызвать сбор всех инфы
		ztc := e.getDataNewCard() // todo: отладить
		fmt.Println(ztc.tags)
	})
	bottom := container.NewBorder(nil, nil, nil, CreateButton, e.Date)

	return container.NewBorder(top, bottom, nil, nil, e.Text)
}

func (e *elmFormType) textForm() *fyne.Container {

	e.Date = newFormatEntry()
	e.Date.SetText(fmt.Sprintf("%v", selectedFile.data))
	e.Title = newFormatEntry()
	e.Text = newText()
	// e.Text.SetText("<Текст>")

	return container.NewBorder(
		e.Title,
		e.Date,
		nil,
		nil,
		e.Text)
}

func (e *elmFormType) addInfoForm() *fyne.Container {

	e.BindNumbers = newNumericalEntry()
	e.Tags = widget.NewMultiLineEntry()
	e.Binds = widget.NewMultiLineEntry()

	binsNumber := container.NewBorder(nil, nil, widget.NewLabel("Номера связных карт:"), nil, e.BindNumbers)

	ent1 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Теги:"), nil, nil, nil, e.Tags))

	ent2 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Связное:"), nil, nil, nil, e.Binds))

	return container.NewBorder(binsNumber, nil, nil, nil, container.NewGridWithColumns(1, ent2, ent1))
}

func (e *elmFormType) sourceForm() *fyne.Container {
	e.SourceNumber = newNumericalEntry()
	e.Source = newText()
	e.Quotation = newText()

	sourceNumber := container.NewBorder(nil, nil, widget.NewLabel("Номер карты источника:"), nil, e.SourceNumber)

	ent1 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, e.Source))

	ent2 := container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(widget.NewLabel("Цитата:"), nil, nil, nil, e.Quotation))

	return container.NewBorder(sourceNumber, nil, nil, nil, container.NewGridWithColumns(1, ent1, ent2))
}

func (e *elmFormType) commentForm() *fyne.Container {
	e.Comment = newText()

	return container.NewGridWithColumns(1,
		container.NewBorder(widget.NewLabel("Комментарий:"), nil, nil, nil, e.Comment))
}

func (e *elmFormType) refreshTabs(z ztcBasicsType) {

	e.Date.SetText(fmt.Sprintf("%v", z.data.Format("2006-01-02 15:04")))
	e.Title.SetText(z.title)
	e.Tags.SetText(sliceInColumn(z.tags))
	e.Binds.SetText(sliceInColumn(z.bind))
	e.BindNumbers.SetText(sliceInString(z.bindNumbers))
	e.SourceNumber.SetText(sliceInString(z.sourceNumber))
	e.Source.SetText(sliceInColumn(z.source))
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
