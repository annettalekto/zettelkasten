package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type elmFormType struct {
	Name      *widget.Entry
	Tags      *widget.Entry
	Binds     *widget.Entry
	Source    *widget.Entry
	Text      *widget.Entry
	Date      *widget.Label
	Quotation *widget.Entry
	Comment   *widget.Entry
}

var elmForm elmFormType

func (e *elmFormType) addInfoForm() *fyne.Container {

	// todo: номер карты добавить
	e.Date = newFormatLabel(fmt.Sprintf("%v", selectedFile.data))

	e.Name = newFormatEntry()
	e.Name.SetText("<Имя файла>")

	e.Text = newText()
	e.Text.SetText("<Текст>")

	btnEdit := widget.NewButton("Редакт.", func() {
		// вызвать окно редактора todo:
		e.editForm()
	})

	bottom := container.NewBorder(nil, nil, e.Date, btnEdit)

	return container.NewBorder(e.Name, bottom, nil, nil, e.Text)
}

func (e *elmFormType) viewForm() *fyne.Container {

	e.Name = newFormatEntry()
	e.Name.SetText("<Имя файла>")

	e.Tags = widget.NewMultiLineEntry()
	tagBox := container.NewBorder(widget.NewLabel("Теги:"), nil, nil, nil, e.Tags)

	e.Binds = widget.NewMultiLineEntry()
	dindsBox := container.NewBorder(widget.NewLabel("Связное:"), nil, nil, nil, e.Binds)

	e.Source = widget.NewMultiLineEntry()
	sourceBox := container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, e.Source)

	box := container.NewGridWithColumns(1, tagBox, dindsBox, sourceBox)

	return container.NewBorder(e.Name, nil, nil, nil, box)
}

func (e *elmFormType) sourceForm() *fyne.Container {
	e.Source = newText()
	sourceBox := container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, e.Source)

	e.Quotation = newText()
	quotationBox := container.NewBorder(widget.NewLabel("Цитата:"), nil, nil, nil, e.Quotation)

	return container.NewGridWithColumns(1, sourceBox, quotationBox)
}

func (e *elmFormType) commentForm() *fyne.Container {
	e.Comment = newText()
	e.Comment.SetText("<Комментарий>")

	return container.NewGridWithColumns(1, e.Comment)
}

func (e *elmFormType) refreshTabs(z ztcBasicsType) {

	e.Date.SetText(fmt.Sprintf("%v", z.data.Format("2006-01-02 15:04"))) //d.Format("2006-01-02 15:04")
	if len(z.title) > 1 {
		e.Name.SetText(z.title)
	}

	if len(z.tags) > 0 {
		e.Tags.SetText(formatSlice(z.tags))
	}

	if len(z.bind) > 0 {
		e.Binds.SetText(formatSlice(z.bind))
	}

	if len(z.source) > 0 {
		e.Source.SetText(formatSlice(z.source))
	}

	e.Text.SetText(getTextFromFile(z.filePath))

	quotation := getQuotationFromFile(z.filePath)
	if len(quotation) > 0 {
		e.Quotation.SetText(quotation)
	}

	comment := getCommentFromFile(z.filePath)
	if len(comment) > 0 {
		e.Comment.SetText(getCommentFromFile(z.filePath))
	}
}

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

	b := container.NewGridWithColumns(1, tagBox, dindsBox, sourceBox, quotationBox, commentBox)
	box := container.NewBorder(e.Name, nil, nil, nil, b) // + текст
	w.SetContent(box)

	w.Show()
}
