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

type addInfoType struct {
	Name   *widget.Entry
	Tags   *widget.Entry
	Binds  *widget.Entry
	Source *widget.Entry
}

type sourceInfoType struct {
	Source    *widget.Entry
	Quotation *widget.Entry
}

type commentType struct {
	Comment *widget.Entry
}

var addInfoForm addInfoType
var viewForm viewType
var sourceInfoForm sourceInfoType
var commentForm commentType

func (v *viewType) initForm() *fyne.Container {

	// номер карты в шапке?
	v.Date = newFormatLabel(fmt.Sprintf("%v", selectedFile.data))
	// d := container.NewBorder(nil, nil, nil, v.Date)

	v.Name = newFormatEntry()
	v.Name.SetText("<Имя файла>")
	// top := container.NewVBox(d, v.Name)

	v.Text = newText()
	v.Text.SetText("<Текст>")

	bottom := container.NewBorder(nil, nil, v.Date, widget.NewButton("Редакт.", nil))

	return container.NewBorder(v.Name, bottom, nil, nil, v.Text)
}

func (a *addInfoType) initForm() *fyne.Container {

	a.Name = newFormatEntry()
	a.Name.SetText("<Имя файла>")

	a.Tags = widget.NewMultiLineEntry()
	tagsBox := container.NewBorder(widget.NewLabel("Теги:"), nil, nil, nil, a.Tags)

	a.Binds = widget.NewMultiLineEntry()
	dindsBox := container.NewBorder(widget.NewLabel("Связное:"), nil, nil, nil, a.Binds)

	a.Source = widget.NewMultiLineEntry()
	sourceBox := container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, a.Source)

	box := container.NewGridWithColumns(1, tagsBox, dindsBox, sourceBox)

	return container.NewBorder(a.Name, nil, nil, nil, box)
}

func (s *sourceInfoType) initForm() *fyne.Container {
	s.Source = newText()
	sourceBox := container.NewBorder(widget.NewLabel("Источники:"), nil, nil, nil, s.Source)

	s.Quotation = newText()
	quotationBox := container.NewBorder(widget.NewLabel("Цитата:"), nil, nil, nil, s.Quotation)

	return container.NewGridWithColumns(1, sourceBox, quotationBox)
}

func (c *commentType) initForm() *fyne.Container {
	c.Comment = newText()
	c.Comment.SetText("<Комментарий>")

	return container.NewGridWithColumns(1, c.Comment)
}

func refreshTabs(z ztcBasicsType) {

	viewForm.Date.SetText(fmt.Sprintf("%v", z.data.Format("2006-01-02 15:04"))) //d.Format("2006-01-02 15:04")
	if len(z.title) > 1 {
		viewForm.Name.SetText(z.title)
		addInfoForm.Name.SetText(z.title)
	}

	if len(z.tags) > 0 {
		addInfoForm.Tags.SetText(formatSlice(z.tags))
	}

	if len(z.bind) > 0 {
		addInfoForm.Binds.SetText(formatSlice(z.bind))
	}

	if len(z.source) > 0 {
		addInfoForm.Source.SetText(formatSlice(z.source))
		sourceInfoForm.Source.SetText(formatSlice(z.source))
	}

	viewForm.Text.SetText(getTextFromFile(z.filePath))

	quotation := getQuotationFromFile(z.filePath)
	if len(quotation) > 0 {
		sourceInfoForm.Quotation.SetText(quotation)
	}

	comment := getCommentFromFile(z.filePath)
	if len(comment) > 0 {
		commentForm.Comment.SetText(getCommentFromFile(z.filePath))
	}
}
